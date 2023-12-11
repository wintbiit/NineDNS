package main

import (
	"github.com/miekg/dns"
	"github.com/wintbiit/ninedns/server"
	"github.com/wintbiit/ninedns/utils"
	"go.uber.org/zap"
)

var servers = make(map[string]*server.Server)

func main() {
	for name, domain := range utils.C.Domains {
		serv, err := server.NewServer(&domain, name)
		if err != nil {
			zap.S().Errorf("Failed to create server for domain %s: %s", name, err)
			continue
		}

		servers[name] = serv
	}

	defer func() {
		for _, serv := range servers {
			serv.Close()
		}
	}()

	zap.S().Infof("Nine DNS started on %s", utils.C.Addr)

	serv := &dns.Server{Addr: utils.C.Addr, Net: "udp"}
	if err := serv.ListenAndServe(); err != nil {
		zap.S().Fatal(err)
	}
}
