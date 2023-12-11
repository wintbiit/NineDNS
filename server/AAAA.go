package server

import "github.com/miekg/dns"

func (s *RuleSet) handleAAAA(r *dns.Msg, q *dns.Question, m *dns.Msg) {
	if s.Recursion {
		s.l.Debugf("Recursion enabled, forwarding request to upstream: %s", s.Upstream)
		resp, _, err := s.dnsClient.Exchange(r, s.Upstream)
		if err != nil {
			s.l.Errorf("Failed to forward request to upstream: %s", err)
			return
		}
		m.Answer = append(m.Answer, resp.Answer...)
	} else {
		record := s.findRecord(q.Name, q.Qtype)

		if record == nil {
			// Record not found
			s.l.Infof("No record found for question: %+v", q)
			m.SetRcode(r, dns.RcodeNameError)
			return
		}

		m.Answer = append(m.Answer, &dns.AAAA{
			Hdr:  s.Header(record),
			AAAA: record.Value.IP(),
		})
	}
}
