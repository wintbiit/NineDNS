package server

import (
	"github.com/miekg/dns"
)

func (s *RuleSet) handleSOA(r *dns.Msg, q *dns.Question, m *dns.Msg) {
	record := s.findRecord(q.Name, q.Qtype)

	if record == nil {
		m.SetRcode(r, dns.RcodeNameError)
		return
	}

	soa, err := record.Value.SOA()
	if err != nil {
		s.l.Errorf("Failed to parse SOA record: %s", err)
		m.SetRcode(r, dns.RcodeServerFailure)
		return
	}
	m.Answer = append(m.Answer, &dns.SOA{
		Hdr:     s.Header(record),
		Ns:      soa.NS,
		Mbox:    soa.MBox,
		Serial:  soa.Serial,
		Refresh: soa.Refresh,
		Retry:   soa.Retry,
		Expire:  soa.Expire,
		Minttl:  soa.MinTTL,
	})
}
