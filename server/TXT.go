package server

import "github.com/miekg/dns"

func (s *RuleSet) handleTXT(r, m *dns.Msg, name string) error {
	records := s.findRecords(name, dns.TypeTXT)

	if records == nil {
		if !s.Recursion {
			return nil
		}

		s.l.Debugf("Recursion enabled, forwarding request to upstream: %s", s.Upstream)
		resp, _, err := s.dnsClient.Exchange(r, s.Upstream)
		if err != nil {
			return err
		}
		m.Answer = append(m.Answer, resp.Answer...)

		return nil
	}

	txt := make([]string, len(records))
	for i, record := range records {
		txt[i] = record.Value.String()
	}

	m.Answer = append(m.Answer, &dns.TXT{
		Hdr: s.Header(&records[0]),
		Txt: txt,
	})

	return nil
}
