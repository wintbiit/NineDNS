package resolver

import (
	"github.com/miekg/dns"
	"github.com/wintbiit/ninedns/model"
)

type CNAME struct{}

func init() {
	resolvers[dns.TypeCNAME] = &CNAME{}
}

func (_ *CNAME) Resolve(s model.RecordProvider, name string) ([]dns.RR, error) {
	record := s.FindRecord(name, dns.TypeCNAME)
	if record == nil {
		return nil, nil
	}

	cname := record.Value.String()

	rr := &dns.CNAME{
		Hdr:    s.Header(record),
		Target: dns.Fqdn(cname),
	}

	return []dns.RR{rr}, nil
}
