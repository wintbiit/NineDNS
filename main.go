package main

import (
	"github.com/miekg/dns"
	"github.com/wintbiit/ninedns/server"
	"github.com/wintbiit/ninedns/utils"
	"go.uber.org/zap"
)

var servers = make(map[string]*server.Server)

func main() {
	for _, domain := range utils.C.Domains {
		domain := domain
		serv, err := server.NewServer(&domain)
		if err != nil {
			zap.S().Errorf("Failed to create server for domain %s: %s", domain.Domain, err)
			continue
		}

		servers[domain.Domain] = serv
	}

	serv := &dns.Server{Addr: utils.C.Addr, Net: "udp"}
	if err := serv.ListenAndServe(); err != nil {
		zap.S().Fatal(err)
	}
}
