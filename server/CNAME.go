package server

import (
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

func (s *RuleSet) handleCNAME(r, m *dns.Msg, name string) error {
	record := s.findRecord(name, dns.TypeCNAME)
	if record == nil {
		if !s.Recursion {
			return fmt.Errorf("no record found for question: %+v", name)
		}

		s.l.Debugf("Recursion enabled, forwarding request to upstream: %s", s.Upstream)
		resp, _, err := s.dnsClient.Exchange(r, s.Upstream)
		if err != nil {
			s.l.Errorf("Failed to forward request to upstream: %s", err)
			return err
		}
		m.Answer = append(m.Answer, resp.Answer...)

		return nil
	}

	s.l.Infof("Response for question: %+v", name)

	cname := record.Value.String()
	cname = dns.Fqdn(cname)

	m.Answer = append(m.Answer, &dns.CNAME{
		Hdr:    s.Header(record),
		Target: cname,
	})

	return nil
}

func (s *RuleSet) localCNAME(name string) string {
	record := s.findRecord(name, dns.TypeCNAME)
	if record == nil {
		return ""
	}

	cname := record.Value.String()
	if !strings.HasSuffix(cname, ".") {
		cname = cname + "." + s.DomainName
	}

	return cname
}
