package server

import (
	"fmt"

	"github.com/miekg/dns"
)

func (s *RuleSet) handleSOA(r, m *dns.Msg, name string) error {
	record := s.findRecord(name, dns.TypeSOA)

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

	soa, err := record.Value.SOA()
	if err != nil {
		return err
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

	return nil
}
