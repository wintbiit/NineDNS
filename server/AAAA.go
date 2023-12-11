package server

import (
	"fmt"

	"github.com/miekg/dns"
)

func (s *RuleSet) handleAAAA(r *dns.Msg, q *dns.Question, m *dns.Msg) error {
	record := s.findRecord(q.Name, q.Qtype)
	if record == nil {
		if !s.Recursion {
			return fmt.Errorf("no record found for question: %+v", q)
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

	s.l.Infof("Response for question: %+v", q)

	m.Answer = append(m.Answer, &dns.AAAA{
		Hdr:  s.Header(record),
		AAAA: record.Value.IP(),
	})

	return nil
}
