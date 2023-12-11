package server

import (
	"fmt"

	"github.com/miekg/dns"
)

type ResolverSOA struct{}

func init() {
	resolvers[dns.TypeSOA] = &ResolverSOA{}
}

func (_ *ResolverSOA) Resolve(s *RuleSet, r *dns.Msg, name string) ([]dns.RR, error) {
	record := s.findRecord(name, dns.TypeSOA)

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

	soa, err := record.Value.SOA()
	if err != nil {
		return nil, err
	}

	rr := &dns.SOA{
		Hdr:     s.Header(record),
		Ns:      dns.Fqdn(soa.NS),
		Mbox:    dns.Fqdn(soa.MBox),
		Serial:  soa.Serial,
		Refresh: soa.Refresh,
		Retry:   soa.Retry,
		Expire:  soa.Expire,
		Minttl:  soa.MinTTL,
	}

	return []dns.RR{rr}, nil
}
