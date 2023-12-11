package server

import "github.com/miekg/dns"

type ResolverTXT struct{}

func init() {
	resolvers[dns.TypeTXT] = &ResolverTXT{}
}

func (_ *ResolverTXT) Resolve(s *RuleSet, r *dns.Msg, name string) ([]dns.RR, error) {
	records := s.findRecords(name, dns.TypeTXT)

	if records == nil {
		if !s.Recursion {
			return nil, nil
		}

		s.l.Debugf("Recursion enabled, forwarding request to upstream: %s", s.Upstream)
		resp, _, err := s.dnsClient.Exchange(r, s.Upstream)
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
