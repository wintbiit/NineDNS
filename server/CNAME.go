package server

import (
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

type ResolverCNAME struct{}

func init() {
	resolvers[dns.TypeCNAME] = &ResolverCNAME{}
}

func (_ *ResolverCNAME) Resolve(s *RuleSet, r *dns.Msg, name string) ([]dns.RR, error) {
	record := s.findRecord(name, dns.TypeCNAME)
	if record == nil {
		if !s.Recursion {
			return nil, fmt.Errorf("no record found for question: %+v", name)
		}

		s.l.Debugf("Recursion enabled, forwarding request to upstream: %s", s.Upstream)
		resp, _, err := s.dnsClient.Exchange(r, s.Upstream)
		if err != nil {
			s.l.Errorf("Failed to forward request to upstream: %s", err)
			return nil, err
		}

		return resp.Answer, nil
	}

	s.l.Infof("Response for question: %+v", name)

	cname := record.Value.String()
	cname = dns.Fqdn(cname)

	rr := &dns.CNAME{
		Hdr:    s.Header(record),
		Target: cname,
	}

	return []dns.RR{rr}, nil
}

func (s *RuleSet) localCNAME(name string) string {
	record := s.findRecord(name, dns.TypeCNAME)
	if record == nil {
		return ""
	}

	cname := record.Value.String()
	if !strings.HasSuffix(cname, ".") {
		cname = cname + "." + s.DomainName
	}

	return cname
}
