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

	txt := make([]string, len(records))
	for i, record := range records {
		txt[i] = record.Value.String()
	}

	rr := &dns.TXT{
		Hdr: s.Header(&records[0]),
		Txt: txt,
	}

	return []dns.RR{rr}, nil
}
