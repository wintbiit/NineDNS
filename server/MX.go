package server

import (
	"fmt"

	"github.com/miekg/dns"
)

func (s *RuleSet) handleMX(r, m *dns.Msg, name string) error {
	record := s.findRecord(name, dns.TypeMX)
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

	mx, err := record.Value.MX()
	if err != nil {
		return err
	}

	m.Answer = append(m.Answer, &dns.MX{
		Hdr:        s.Header(record),
		Preference: mx.Preference,
		Mx:         mx.MX,
	})

	return nil
}
