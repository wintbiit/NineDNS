package server

import (
	"fmt"

	"github.com/miekg/dns"
)

func (s *RuleSet) handleSRV(r *dns.Msg, q *dns.Question, m *dns.Msg) error {
	record := s.findRecord(q.Name, q.Qtype)

	if record == nil {
		if !s.Recursion {
			return fmt.Errorf("no record found for question: %+v", q)
		}

		s.l.Debugf("Recursion enabled, forwarding request to upstream: %s", s.Upstream)
		resp, _, err := s.dnsClient.Exchange(r, s.Upstream)
		if err != nil {
			return err
		}
		m.Answer = append(m.Answer, resp.Answer...)

		return nil
	}

	srv, err := record.Value.SRV()
	if err != nil {
		return err
	}

	m.Answer = append(m.Answer, &dns.SRV{
		Hdr:      s.Header(record),
		Priority: srv.Priority,
		Weight:   srv.Weight,
		Port:     srv.Port,
		Target:   srv.Target,
	})

	return nil
}
