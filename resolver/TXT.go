package resolver

import (
	"github.com/miekg/dns"
	"github.com/wintbiit/ninedns/model"
)

type TXT struct{}

func init() {
	resolvers[dns.TypeTXT] = &TXT{}
}

func (_ *TXT) Resolve(s model.RecordProvider, r *dns.Msg, name string) ([]dns.RR, error) {
	records := s.FindRecords(name, dns.TypeTXT)

	if records == nil {
		if !s.Recursion() {
			return nil, nil
		}

		resp, err := s.Exchange(r)
		if err != nil {
			return nil, err
		}

		return resp.Answer, nil
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
