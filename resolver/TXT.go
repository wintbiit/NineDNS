package resolver

import (
	"github.com/miekg/dns"
	"github.com/wintbiit/ninedns/model"
)

type TXT struct{}

func init() {
	resolvers[dns.TypeTXT] = &TXT{}
}

func (_ *TXT) Resolve(s model.RecordProvider, name string) ([]dns.RR, error) {
	records := s.FindRecords(name, dns.TypeTXT)

	if records == nil {
		return nil, nil
	}

	rrs := make([]dns.RR, len(records))
	for i := range records {
		rrs[i] = &dns.TXT{
			Hdr: s.Header(&records[i]),
			Txt: []string{records[i].Value.String()},
		}
	}

	return rrs, nil
}
