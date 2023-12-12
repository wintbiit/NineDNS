package resolver

import (
	"github.com/miekg/dns"
	"github.com/wintbiit/ninedns/model"
)

type SRV struct{}

func init() {
	resolvers[dns.TypeSRV] = &SRV{}
}

func (_ *SRV) Resolve(s model.RecordProvider, name string) ([]dns.RR, error) {
	record := s.FindRecord(name, dns.TypeSRV)

	if record == nil {
		return nil, nil
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
		Target:   dns.Fqdn(srv.Target),
	}

	return []dns.RR{rr}, nil
}
