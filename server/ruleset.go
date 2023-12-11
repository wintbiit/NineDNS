package server

import (
	"fmt"
	"math/rand"
	"net"
	"sort"
	"strings"

	"github.com/miekg/dns"

	"github.com/wintbiit/ninedns/model"
	"go.uber.org/zap"
)

type RuleSet struct {
	model.Rule
	*Server
	Name  string
	l     *zap.SugaredLogger
	cidrs []*net.IPNet
}

type Resolver interface {
	Resolve(*RuleSet, *dns.Msg, string) ([]dns.RR, error)
}

var resolvers = make(map[uint16]Resolver)

func (s *Server) newRuleSet(name string, rule model.Rule) *RuleSet {
	ruleSet := &RuleSet{
		Rule:   rule,
		Server: s,
		Name:   name,
		l:      s.l.Named(s.DomainName + "/" + name),
		cidrs:  make([]*net.IPNet, len(rule.CIDRs)),
	}

	// Init CIDR rules
	for i, cidr := range ruleSet.CIDRs {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			s.l.Errorf("Failed to parse CIDR %s: %s", cidr, err)
			continue
		}

		ruleSet.cidrs[i] = ipNet
	}

	return ruleSet
}

func (s *RuleSet) query(r, m *dns.Msg) {
	for _, q := range r.Question {
		// 1. Try CNAME
		name := q.Name
		cname := s.localCNAME(q.Name)
		if cname != "" {
			s.l.Infof("Question %s CNAME hit: %s", q.String(), cname)
			q.Name = cname
		}

		err := s.question(&q, r, m, name)
		if err != nil {
			s.l.Warnf("Failed to handle question %s: %s", q.String(), err)
			m.SetRcode(r, dns.RcodeNameError)
			continue
		}
	}
}

func (s *RuleSet) question(q *dns.Question, r, m *dns.Msg, name string) error {
	resolver, ok := resolvers[q.Qtype]
	if !ok {
		return fmt.Errorf("unsupported DNS question type: %s", dns.TypeToString[q.Qtype])
	}

	records, err := resolver.Resolve(s, r, q.Name)
	if err != nil {
		return err
	}

	for _, record := range records {
		record.Header().Name = name
		m.Answer = append(m.Answer, record)
	}

	return nil
}

func (s *RuleSet) ShouldHandle(ip net.IP, port int, zone, network string) bool {
	matchers := 0
	matched := 0

	if len(s.cidrs) > 0 {
		matchers++
		for _, cidr := range s.cidrs {
			if cidr.Contains(ip) {
				matched++
				break
			}
		}
	}

	if len(s.Ports) > 0 {
		matchers++
		for _, portRule := range s.Ports {
			if portRule.Contains(port) {
				matched++
				break
			}
		}
	}

	if len(s.Types) > 0 {
		matchers++
		for _, typ := range s.Types {
			if typ == network {
				matched++
				break
			}
		}
	}

	return matchers == matched
}

func (s *RuleSet) findRecords(name string, quesType uint16) []model.Record {
	name = strings.TrimSuffix(name, s.DomainName)
	name = strings.TrimSuffix(name, ".")
	if name == "" {
		name = "@"
	}
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

	for _, prov := range s.providers {
		recs, err := prov.Provide(s.Name)
		if err != nil {
			s.l.Errorf("Failed to provide records: %s", err)
			continue
		}

		records = append(records, recs...)
	}

	for index, record := range records {
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

		records[index] = record
	}

	s.l.Infof("Refreshed %d records", len(records))
	s.l.Debug("Refreshed records: %v", records)
}
