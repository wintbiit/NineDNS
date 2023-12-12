package resolver

import (
	"github.com/miekg/dns"
	"github.com/wintbiit/ninedns/model"
)

type MX struct{}

func init() {
	resolvers[dns.TypeMX] = &MX{}
}

func (_ *MX) Resolve(s model.RecordProvider, name string) ([]dns.RR, error) {
	record := s.FindRecord(name, dns.TypeMX)
	if record == nil {
		return nil, nil
	}

	mx, err := record.Value.MX()
	if err != nil {
		return nil, err
	}

	rr := &dns.MX{
		Hdr:        s.Header(record),
		Preference: mx.Preference,
		Mx:         dns.Fqdn(mx.MX),
	}

	return []dns.RR{rr}, nil
}
