package resolver

import (
	"github.com/miekg/dns"
	"github.com/wintbiit/ninedns/model"
)

type Resolver interface {
	Resolve(model.RecordProvider, *dns.Msg, string) ([]dns.RR, error)
}

var resolvers = make(map[uint16]Resolver)

func Resolve(typ uint16, p model.RecordProvider, r *dns.Msg, name string) ([]dns.RR, error) {
	resolver, ok := resolvers[typ]
	if !ok {
		return nil, nil
	}

	return resolver.Resolve(p, r, name)
}
