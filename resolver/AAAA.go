package resolver

import (
	"fmt"

	"github.com/miekg/dns"
	"github.com/wintbiit/ninedns/model"
)

type AAAA struct{}

func init() {
	resolvers[dns.TypeAAAA] = &AAAA{}
}

func (_ *AAAA) Resolve(s model.RecordProvider, r *dns.Msg, name string) ([]dns.RR, error) {
	record := s.FindRecord(name, dns.TypeAAAA)
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

	rr := &dns.AAAA{
		Hdr:  s.Header(record),
		AAAA: record.Value.IP(),
	}

	return []dns.RR{rr}, nil
}
