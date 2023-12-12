package resolver

import (
	"github.com/miekg/dns"
	"github.com/wintbiit/ninedns/model"
)

type SOA struct{}

func init() {
	resolvers[dns.TypeSOA] = &SOA{}
}

func (_ *SOA) Resolve(s model.RecordProvider, name string) ([]dns.RR, error) {
	record := s.FindRecord(name, dns.TypeSOA)

	if record == nil {
		return nil, nil
	}

	soa, err := record.Value.SOA()
	if err != nil {
		return nil, err
	}

	rr := &dns.SOA{
		Hdr:     s.Header(record),
		Ns:      dns.Fqdn(soa.NS),
		Mbox:    dns.Fqdn(soa.MBox),
		Serial:  soa.Serial,
		Refresh: soa.Refresh,
		Retry:   soa.Retry,
		Expire:  soa.Expire,
		Minttl:  soa.MinTTL,
	}

	return []dns.RR{rr}, nil
}
