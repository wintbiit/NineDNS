package server

import (
	"fmt"

	"github.com/miekg/dns"
)

type ResolverMX struct{}

func init() {
	resolvers[dns.TypeMX] = &ResolverMX{}
}

func (_ *ResolverMX) Resolve(s *RuleSet, r *dns.Msg, name string) ([]dns.RR, error) {
	record := s.findRecord(name, dns.TypeMX)
	if record == nil {
		if !s.Recursion {
			return nil, fmt.Errorf("no record found for question: %+v", name)
		}

		s.l.Debugf("Recursion enabled, forwarding request to upstream: %s", s.Upstream)
		resp, _, err := s.dnsClient.Exchange(r, s.Upstream)
		if err != nil {
			s.l.Errorf("Failed to forward request to upstream: %s", err)
			return nil, err
		}

		return resp.Answer, nil
	}

	s.l.Infof("Response for question: %+v", name)

	mx, err := record.Value.MX()
	if err != nil {
		return nil, err
	}

	rr := &dns.MX{
		Hdr:        s.Header(record),
		Preference: mx.Preference,
		Mx:         mx.MX,
	}

	return []dns.RR{rr}, nil
}
