package server

import (
	"net"
	"strings"
	"time"

	"github.com/wintbiit/ninedns/provider"

	"github.com/redis/go-redis/v9"

	"github.com/miekg/dns"
	"github.com/wintbiit/ninedns/cache"
	"github.com/wintbiit/ninedns/log"
	"github.com/wintbiit/ninedns/model"
	"go.uber.org/zap"
)

type Server struct {
	model.Domain
	DomainName  string
	l           *zap.SugaredLogger
	dnsClient   *dns.Client
	cacheClient cache.API
	providers   map[string]provider.Provider
	rules       map[string]*RuleSet
}

func NewServer(config *model.Domain, domain string) (*Server, error) {
	server := &Server{
		Domain:     *config,
		DomainName: domain,
		l:          log.NewLogger(domain).Sugar(),
		rules:      make(map[string]*RuleSet),
		dnsClient:  new(dns.Client),
		providers:  make(map[string]provider.Provider),
	}

	server.checkConfig()

	// Init Providers
	for name, conf := range server.Domain.Providers {
		prov, err := provider.NewProvider(name, conf)
		if err != nil {
			server.l.Errorf("Failed to create provider %s: %s", name, err)
			return nil, err
		}

		server.providers[name] = prov
	}

	// Init Cache
	cacheClient, err := cache.NewClient(server.DomainName, server.TTL)
	if err != nil {
		server.l.Errorf("Failed to create cache client: %s", err)
		return nil, err
	}
	server.cacheClient = cacheClient

	// Init Rules
	for name, rule := range server.Domain.Rules {
		set := server.newRuleSet(name, rule)

		go func() {
			set.RefreshRecords()
			for range time.Tick(time.Duration(server.TTL) * time.Second) {
				set.RefreshRecords()
			}
		}()

		server.rules[name] = set
	}

	dns.HandleFunc(server.DomainName, server.handle)

	return server, nil
}

func (s *Server) checkConfig() {
	if !strings.HasSuffix(s.DomainName, ".") {
		s.l.Warn("Record domain missing `.` suffix, automatically add it.")
		s.DomainName += "."
	}

	if s.Domain.Authoritative {
		s.l.Warn("Server is authoritative, please ensure it's correct.")
	}

	if s.Domain.Recursion {
		s.l.Warn("Server is recursion, please ensure it's correct.")
	}

	if s.Domain.TTL == 0 {
		s.Domain.TTL = 60
		s.l.Warn("Server TTL is 0, automatically set it to 60.")
	}

	if s.Domain.Upstream == "" {
		s.Domain.Upstream = "223.5.5.5:53"
		s.l.Warn("Server upstream is empty, automatically set it to %s.", s.Domain.Upstream)
	}

	if s.Domain.Rules == nil {
		s.l.Warn("Server rules is empty, please ensure it's correct.")
	}
}

func (s *Server) handle(w dns.ResponseWriter, r *dns.Msg) {
	remoteAddr := w.RemoteAddr()
	s.l.Debugf("Receive DNS request {%+v} from %s: %s", r, remoteAddr.Network(), remoteAddr.String())

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = s.Authoritative
	m.RecursionAvailable = s.Recursion

	handler := s.MatchHandler(w)

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

func (s *Server) Header(r *model.Record) dns.RR_Header {
	return dns.RR_Header{
		Name:   r.Host + "." + s.DomainName,
		Rrtype: r.Type.DnsType(),
		Class:  dns.ClassINET,
		Ttl:    s.TTL,
	}
}

func (s *Server) MatchHandler(w dns.ResponseWriter) *RuleSet {
	addr := w.RemoteAddr()
	var ip net.IP
	var port int
	var zone, network string
	if addr.Network() == "tcp" {
		addr := addr.(*net.TCPAddr)
		ip = addr.IP
		port = addr.Port
		zone = addr.Zone
		network = addr.Network()
	} else if addr.Network() == "udp" {
		addr := addr.(*net.UDPAddr)
		ip = addr.IP
		port = addr.Port
		zone = addr.Zone
	} else {
		s.l.Warnf("Unknown network type %s", addr.Network())
		return nil
	}

	var handlerName string
	var err error

	handlerName, err = s.cacheClient.GetRuntimeCache("handler:" + ip.String())
	if err != nil {
		if err == redis.Nil {
			handlerName = s.matchHandler(ip, port, zone, network)
			if handlerName == "" {
				return nil
			}
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

func (s *Server) matchHandler(ip net.IP, port int, zone, network string) string {
	ruleName := ""
	for _, rule := range s.rules {
		if rule.ShouldHandle(ip, port, zone, network) {
			ruleName = rule.Name
			break
		}
	}

	return ruleName
}

func (s *Server) Close() error {
	return s.cacheClient.Close()
}
