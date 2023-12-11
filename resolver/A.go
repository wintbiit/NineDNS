package resolver

import (
	"fmt"

	"github.com/miekg/dns"
	"github.com/wintbiit/ninedns/model"
)

type A struct{}

func init() {
	resolvers[dns.TypeA] = &A{}
}

func (_ *A) Resolve(s model.RuleProvider, r *dns.Msg, name string) ([]dns.RR, error) {
	record := s.FindRecord(name, dns.TypeA)
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

	rr := &dns.A{
		Hdr: s.Header(record),
		A:   record.Value.IP(),
	}

	return []dns.RR{rr}, nil
}
