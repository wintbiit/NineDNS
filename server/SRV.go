package server

import (
	"fmt"

	"github.com/miekg/dns"
)

type ResolverSRV struct{}

func init() {
	resolvers[dns.TypeSRV] = &ResolverSRV{}
}

func (_ *ResolverSRV) Resolve(s *RuleSet, r *dns.Msg, name string) ([]dns.RR, error) {
	record := s.findRecord(name, dns.TypeSRV)

	if record == nil {
		if !s.Recursion {
			return nil, fmt.Errorf("no record found for question: %+v", name)
		}

		s.l.Debugf("Recursion enabled, forwarding request to upstream: %s", s.Upstream)
		resp, _, err := s.dnsClient.Exchange(r, s.Upstream)
		if err != nil {
			return nil, err
		}

		return resp.Answer, nil
	}

	srv, err := record.Value.SRV()
	if err != nil {
		return nil, err
	}

	rr := &dns.SRV{
		Hdr:      s.Header(record),
		Priority: srv.Priority,
		Weight:   srv.Weight,
		Port:     srv.Port,
		Target:   srv.Target,
	}

	return []dns.RR{rr}, nil
}
