package server

import (
	"math/rand"
	"net"
	"sort"
	"strings"

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
	l           *zap.SugaredLogger
	dnsClient   *dns.Client
	cacheClient cache.API
	dbClient    *gorm.DB
	rules       map[string]*RuleSet
}

type RuleSet struct {
	model.Rule
	*Server
	cidr *net.IPNet
}

func NewServer(config *model.Domain) (*Server, error) {
	server := &Server{
		Domain:    *config,
		l:         log.NewLogger(config.Domain).Sugar(),
		dnsClient: new(dns.Client),
	}
	if !strings.HasSuffix(server.Domain.Domain, ".") {
		server.l.Warn("Record domain missing `.` suffix, automatically add it.")
		server.Domain.Domain += "."
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

	for _, rule := range server.Domain.Rules {
		var set RuleSet

		if set.CIDR == "" {
			set.CIDR = "0.0.0.0/0"
			server.l.Warn("Rule CIDR is empty, automatically set it to %s.", set.CIDR)
		}

		_, cidr, err := net.ParseCIDR(rule.CIDR)
		if err != nil {
			server.l.Errorf("Failed to parse CIDR %s: %s", rule.CIDR, err)
			return nil, err
		}

		if err := server.dbClient.Table(rule.Name).AutoMigrate(&model.Record{}); err != nil {
			server.l.Errorf("Failed to auto migrate record: %s", err)
			return nil, err
		}

		set.cidr = cidr

		set.Rule = rule
		set.Server = server
		go set.RefreshRecords()

		server.rules[rule.CIDR] = &set
	}

	cacheClient, err := cache.NewClient(server.Domain.Domain)
	if err != nil {
		server.l.Errorf("Failed to create cache client: %s", err)
		return nil, err
	}

	server.cacheClient = cacheClient

	dns.HandleFunc(server.Domain.Domain, server.handle)

	return server, nil
}

func (s *Server) handle(w dns.ResponseWriter, r *dns.Msg) {
	remoteAddr := w.RemoteAddr()
	s.l.Debugf("Receive DNS request: %+v", r)
	s.l.Debugf("Receive DNS request from %s", remoteAddr)

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = s.Authoritative
	m.RecursionAvailable = s.Recursion

	var handler *RuleSet
	for _, rule := range s.rules {
		if rule.cidr.Contains(net.ParseIP(remoteAddr.String())) {
			handler = rule
			break
		}
	}

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
	records, err := s.cacheClient.FindRecords(name, model.ReadRecordType(quesType).String(), s.CIDR)
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
	s.l.Info("Starting refresh records for rule set %s", s.CIDR)
	var records []model.Record
	defer func() {
		records = append(records, s.Records...)

		for _, record := range records {
			if !record.Enabled {
				continue
			}
			record := record
			if err := s.cacheClient.AddRecord(s.CIDR, &record); err != nil {
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

	if err := tx.Table(s.Name).Where("*").Find(&records).Error; err != nil {
		s.l.Errorf("Failed to query records: %s", err)
		return
	}
}

func (s *Server) Header(r *model.Record) dns.RR_Header {
	return dns.RR_Header{
		Name:   r.Host,
		Rrtype: r.Type.DnsType(),
		Class:  dns.ClassINET,
		Ttl:    s.TTL,
	}
}
