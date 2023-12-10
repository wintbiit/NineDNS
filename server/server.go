package server

import (
	"github.com/redis/go-redis/v9"
	"math/rand"
	"net"
	"sort"
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/wintbiit/ninedns/cache"
	"github.com/wintbiit/ninedns/log"
	"github.com/wintbiit/ninedns/model"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Server struct {
	model.Domain
	DomainName  string
	l           *zap.SugaredLogger
	dnsClient   *dns.Client
	cacheClient cache.API
	dbClient    *gorm.DB
	rules       map[string]*RuleSet
}

type RuleSet struct {
	model.Rule
	*Server
	Name  string
	l     *zap.SugaredLogger
	cidrs []*net.IPNet
}

func NewServer(config *model.Domain, domain string) (*Server, error) {
	server := &Server{
		Domain:     *config,
		DomainName: domain,
		l:          log.NewLogger(domain).Sugar(),
		rules:      make(map[string]*RuleSet),
		dnsClient:  new(dns.Client),
	}
	if !strings.HasSuffix(server.DomainName, ".") {
		server.l.Warn("Record domain missing `.` suffix, automatically add it.")
		server.DomainName += "."
	}

	if server.Domain.Authoritative {
		server.l.Warn("Server is authoritative, please ensure it's correct.")
	}

	if server.Domain.Recursion {
		server.l.Warn("Server is recursion, please ensure it's correct.")
	}

	if server.Domain.TTL == 0 {
		server.Domain.TTL = 60
		server.l.Warn("Server TTL is 0, automatically set it to 60.")
	}

	if server.Domain.Upstream == "" {
		server.Domain.Upstream = "223.5.5.5:53"
		server.l.Warn("Server upstream is empty, automatically set it to %s.", server.Domain.Upstream)
	}

	if server.Domain.Rules == nil {
		server.l.Warn("Server rules is empty, please ensure it's correct.")
	}

	if server.MySQL != "" {
		db, err := gorm.Open(mysql.Open(server.MySQL), &gorm.Config{})
		if err != nil {
			server.l.Errorf("Failed to open MySQL %s: %s", server.MySQL, err)
			return nil, err
		}

		server.dbClient = db
	}

	cacheClient, err := cache.NewClient(server.DomainName, server.TTL)
	if err != nil {
		server.l.Errorf("Failed to create cache client: %s", err)
		return nil, err
	}

	server.cacheClient = cacheClient

	for name, rule := range server.Domain.Rules {
		set := RuleSet{
			Name:   name,
			Server: server,
			l:      server.l.Named(name),
		}

		set.cidrs = make([]*net.IPNet, len(rule.CIDRs))
		for i, cid := range rule.CIDRs {
			if cid == "" {
				cid = "0.0.0.0/0"
				server.l.Warnf("Rule CIDR is empty, automatically set it to %s.", cid)
			}

			_, cidr, err := net.ParseCIDR(cid)
			if err != nil {
				server.l.Errorf("Failed to parse CIDR %s: %s", cidr, err)
				return nil, err
			}

			if err := server.dbClient.Table(name).AutoMigrate(&model.Record{}); err != nil {
				server.l.Errorf("Failed to auto migrate record: %s", err)
				return nil, err
			}

			set.cidrs[i] = cidr
		}

		go func() {
			set.RefreshRecords()
			for range time.Tick(time.Duration(server.TTL) * time.Second) {
				set.RefreshRecords()
			}
		}()

		server.rules[name] = &set
	}

	dns.HandleFunc(server.DomainName, server.handle)

	return server, nil
}

func (s *Server) handle(w dns.ResponseWriter, r *dns.Msg) {
	remoteAddr := w.RemoteAddr()
	s.l.Debugf("Receive DNS request {%+v} from %s: %s", r, remoteAddr.Network(), remoteAddr.String())
	var remoteIp net.IP
	if remoteAddr.Network() == "udp" {
		remoteIp = remoteAddr.(*net.UDPAddr).IP
	} else if remoteAddr.Network() == "tcp" {
		remoteIp = remoteAddr.(*net.TCPAddr).IP
	} else {
		s.l.Warnf("Unsupported network %s", remoteAddr.Network())
		return
	}

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = s.Authoritative
	m.RecursionAvailable = s.Recursion

	handler := s.MatchHandler(remoteIp)

	if handler == nil {
		s.l.Warnf("No rule found for %s", remoteAddr)
		m.SetRcode(r, dns.RcodeRefused)
		return
	}

	s.l.Debugf("Found rule for %s: %+v", remoteAddr, handler)

	for _, q := range r.Question {
		switch q.Qtype {
		case dns.TypeA:
			s.l.Debugf("Receive DNS question type: %s", dns.TypeToString[q.Qtype])
			handler.handleA(r, &q, m)
		case dns.TypeAAAA:
			s.l.Debugf("Receive DNS question type: %s", dns.TypeToString[q.Qtype])
			handler.handleAAAA(r, &q, m)
		case dns.TypeCNAME:
			s.l.Debugf("Receive DNS question type: %s", dns.TypeToString[q.Qtype])
			handler.handleCNAME(r, &q, m)
		case dns.TypeTXT:
			s.l.Debugf("Receive DNS question type: %s", dns.TypeToString[q.Qtype])
			handler.handleTXT(r, &q, m)
		case dns.TypeNS:
			s.l.Debugf("Receive DNS question type: %s", dns.TypeToString[q.Qtype])
			handler.handleNS(r, &q, m)
		case dns.TypeMX:
			s.l.Debugf("Receive DNS question type: %s", dns.TypeToString[q.Qtype])
			handler.handleMX(r, &q, m)
		case dns.TypeSRV:
			s.l.Debugf("Receive DNS question type: %s", dns.TypeToString[q.Qtype])
			handler.handleSRV(r, &q, m)
		case dns.TypeSOA:
			s.l.Debugf("Receive DNS question type: %s", dns.TypeToString[q.Qtype])
			handler.handleSOA(r, &q, m)
		default:
			s.l.Warnf("Unsupported DNS question type: %s", dns.TypeToString[q.Qtype])
		}
	}

	if err := w.WriteMsg(m); err != nil {
		s.l.Errorf("Failed to write response: %s", err)
	}
}

func (s *RuleSet) findRecords(name string, quesType uint16) []model.Record {
	name = strings.TrimSuffix(name, s.DomainName)
	name = strings.TrimSuffix(name, ".")
	records, err := s.cacheClient.FindRecords(name, model.ReadRecordType(quesType).String(), s.Name)
	if err != nil {
		s.l.Errorf("Failed to query records: %s", err)
		return nil
	}

	return records
}

func (s *RuleSet) findRecord(name string, quesType uint16) *model.Record {
	if records := s.findRecords(name, quesType); records != nil {
		if len(records) == 1 {
			return &(records)[0]
		} else if len(records) > 1 {
			var weightSum uint16 = 0
			for _, record := range records {
				weightSum += record.Weight
			}

			sort.SliceStable(records, func(i, j int) bool {
				return (records)[i].Weight > (records)[j].Weight
			})

			random := rand.Intn(int(weightSum))
			for _, record := range records {
				random -= int(record.Weight)
				if random <= 0 {
					return &record
				}
			}
		}
	}

	return nil
}

func (s *RuleSet) RefreshRecords() {
	s.l.Info("Starting refresh records for rule set %s", s.Name)
	var records []model.Record
	defer func() {
		records = append(records, s.Records...)

		for _, record := range records {
			if record.Disabled {
				continue
			}
			if strings.HasSuffix(record.Host, s.DomainName) {
				record.Host = strings.TrimSuffix(record.Host, s.DomainName)
				s.l.Warnf("DNS record %s does not need domain suffix, automatically remove it.", record.Host)
			}

			if strings.HasSuffix(record.Host, ".") {
				record.Host = strings.TrimSuffix(record.Host, ".")
				s.l.Warnf("DNS record %s does not need `.` suffix, automatically remove it.", record.Host)
			}

			if err := s.cacheClient.AddRecord(s.Name, &record); err != nil {
				s.l.Errorf("Failed to add record %s: %s", record.Host, err)
			}
		}

		s.l.Infof("Refreshed %d records", len(records))
	}()

	if s.dbClient == nil {
		return
	}

	tx := s.dbClient.Begin()
	defer tx.Rollback()

	if err := tx.Table(s.Name).Find(&records).Error; err != nil {
		s.l.Errorf("Failed to query records: %s", err)
		return
	}
}

func (s *Server) Header(r *model.Record) dns.RR_Header {
	return dns.RR_Header{
		Name:   r.Host + "." + s.DomainName,
		Rrtype: r.Type.DnsType(),
		Class:  dns.ClassINET,
		Ttl:    s.TTL,
	}
}

func (s *Server) MatchHandler(ip net.IP) *RuleSet {
	finder := func(ip net.IP) string {
		ruleName := ""
		for _, rule := range s.rules {
			for _, cidr := range rule.cidrs {
				if cidr.Contains(ip) {
					ruleName = rule.Name
				}
			}
		}

		return ruleName
	}

	var handlerName string
	var err error

	handlerName, err = s.cacheClient.GetRuntimeCache("handler:" + ip.String())
	if err != nil {
		if err == redis.Nil {
			handlerName = finder(ip)
			if err := s.cacheClient.AddRuntimeCache("handler:"+ip.String(), handlerName, time.Duration(s.TTL)*time.Second); err != nil {
				s.l.Errorf("Failed to add runtime cache: %s", err)
			}
		} else {
			s.l.Errorf("Failed to get runtime cache: %s", err)
		}
	}

	if handlerName == "" {
		return nil
	}

	handler, ok := s.rules[handlerName]
	if !ok {
		return nil
	}

	return handler
}

func (s *Server) Close() error {
	return s.cacheClient.Close()
}
