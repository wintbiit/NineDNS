package resolver

import (
	"github.com/miekg/dns"
	"github.com/wintbiit/ninedns/model"
)

type NS struct{}

func init() {
	resolvers[dns.TypeNS] = &NS{}
}

func (_ *NS) Resolve(s model.RecordProvider, name string) ([]dns.RR, error) {
	record := s.FindRecord(name, dns.TypeNS)
	if record == nil {
		return nil, nil
	}

	rr := &dns.NS{
		Hdr: s.Header(record),
		Ns:  dns.Fqdn(record.Value.String()),
	}

	return []dns.RR{rr}, nil
}
