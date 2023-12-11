package resolver

import (
	"fmt"

	"github.com/miekg/dns"
	"github.com/wintbiit/ninedns/model"
)

type SRV struct{}

func init() {
	resolvers[dns.TypeSRV] = &SRV{}
}

func (_ *SRV) Resolve(s model.RuleProvider, r *dns.Msg, name string) ([]dns.RR, error) {
	record := s.FindRecord(name, dns.TypeSRV)

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
