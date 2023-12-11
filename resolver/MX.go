package resolver

import (
	"fmt"

	"github.com/miekg/dns"
	"github.com/wintbiit/ninedns/model"
)

type MX struct{}

func init() {
	resolvers[dns.TypeMX] = &MX{}
}

func (_ *MX) Resolve(s model.RecordProvider, r *dns.Msg, name string) ([]dns.RR, error) {
	record := s.FindRecord(name, dns.TypeMX)
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

	mx, err := record.Value.MX()
	if err != nil {
		return nil, err
	}

	rr := &dns.MX{
		Hdr:        s.Header(record),
		Preference: mx.Preference,
		Mx:         mx.MX,
	}

	return []dns.RR{rr}, nil
}
