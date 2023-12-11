package server

import (
	"fmt"

	"github.com/miekg/dns"
)

func (s *RuleSet) handleNS(r, m *dns.Msg, name string) error {
	record := s.findRecord(name, dns.TypeNS)
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

	m.Answer = append(m.Answer, &dns.NS{
		Hdr: s.Header(record),
		Ns:  dns.Fqdn(record.Value.String()),
	})

	return nil
}
