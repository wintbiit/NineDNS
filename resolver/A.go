package resolver

import (
	"github.com/miekg/dns"
	"github.com/wintbiit/ninedns/model"
)

type A struct{}

func init() {
	resolvers[dns.TypeA] = &A{}
}

func (_ *A) Resolve(s model.RecordProvider, name string) ([]dns.RR, error) {
	record := s.FindRecord(name, dns.TypeA)
	if record == nil {
		return nil, nil
	}

	rr := &dns.A{
		Hdr: s.Header(record),
		A:   record.Value.IP(),
	}

	return []dns.RR{rr}, nil
}
