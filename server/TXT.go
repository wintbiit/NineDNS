package server

import "github.com/miekg/dns"

func (s *RuleSet) handleTXT(r *dns.Msg, q *dns.Question, m *dns.Msg) {
	records := s.findRecords(q.Name, q.Qtype)

	if records == nil {
		m.SetRcode(r, dns.RcodeNameError)
		return
	}

	txt := make([]string, len(records))
	for i, record := range records {
		txt[i] = record.Value.String()
	}

	m.Answer = append(m.Answer, &dns.TXT{
		Hdr: s.Header(&records[0]),
		Txt: txt,
	})
}
