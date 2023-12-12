package resolver

import (
	"github.com/miekg/dns"
	"github.com/wintbiit/ninedns/model"
)

type AAAA struct{}

func init() {
	resolvers[dns.TypeAAAA] = &AAAA{}
}

func (_ *AAAA) Resolve(s model.RecordProvider, name string) ([]dns.RR, error) {
	record := s.FindRecord(name, dns.TypeAAAA)
	if record == nil {
		return nil, nil
	}

	rr := &dns.AAAA{
		Hdr:  s.Header(record),
		AAAA: record.Value.IP(),
	}

	return []dns.RR{rr}, nil
}
