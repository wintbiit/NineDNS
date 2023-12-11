package server

import (
	"fmt"

	"github.com/miekg/dns"
)

type ResolverA struct{}

func init() {
	resolvers[dns.TypeA] = &ResolverA{}
}

func (_ *ResolverA) Resolve(s *RuleSet, r *dns.Msg, name string) ([]dns.RR, error) {
	record := s.findRecord(name, dns.TypeA)
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

	rr := &dns.A{
		Hdr: s.Header(record),
		A:   record.Value.IP(),
	}

	return []dns.RR{rr}, nil
}
