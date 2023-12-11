package resolver

import (
	"fmt"

	"github.com/miekg/dns"
	"github.com/wintbiit/ninedns/model"
)

type NS struct{}

func init() {
	resolvers[dns.TypeNS] = &NS{}
}

func (_ *NS) Resolve(s model.RuleProvider, r *dns.Msg, name string) ([]dns.RR, error) {
	record := s.FindRecord(name, dns.TypeNS)
	if record == nil {
		if !s.Recursion() {
			return nil, fmt.Errorf("no record found for question: %+v", name)
		}

		resp, err := s.Exchange(r)
		if err != nil {
			return nil, err
		}

		return resp.Answer, nil
	}

	rr := &dns.NS{
		Hdr: s.Header(record),
		Ns:  dns.Fqdn(record.Value.String()),
	}

	return []dns.RR{rr}, nil
}
