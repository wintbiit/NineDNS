package server

import "github.com/miekg/dns"

func (s *RuleSet) handleSRV(r *dns.Msg, q *dns.Question, m *dns.Msg) {
	record := s.findRecord(q.Name, q.Qtype)

	if record == nil {
		m.SetRcode(r, dns.RcodeNameError)
		return
	}

	srv, err := record.Value.SRV()
	if err != nil {
		s.l.Errorf("Failed to parse SRV record: %s", err)
		m.SetRcode(r, dns.RcodeServerFailure)
		return
	}

	m.Answer = append(m.Answer, &dns.SRV{
		Hdr:      s.Header(record),
		Priority: srv.Priority,
		Weight:   srv.Weight,
		Port:     srv.Port,
		Target:   srv.Target,
	})
}
