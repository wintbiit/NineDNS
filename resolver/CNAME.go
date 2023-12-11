package resolver

import (
	"fmt"

	"github.com/miekg/dns"
	"github.com/wintbiit/ninedns/model"
)

type CNAME struct{}

func init() {
	resolvers[dns.TypeCNAME] = &CNAME{}
}

func (_ *CNAME) Resolve(s model.RecordProvider, r *dns.Msg, name string) ([]dns.RR, error) {
	record := s.FindRecord(name, dns.TypeCNAME)
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

	cname := record.Value.String()
	cname = dns.Fqdn(cname)

	rr := &dns.CNAME{
		Hdr:    s.Header(record),
		Target: cname,
	}

	return []dns.RR{rr}, nil
}
