package main

import (
	"os"

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

	pidf, err := os.OpenFile("./.ninedns.pid", os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		zap.S().Fatalf("Failed to open pid file: %s", err)
	}

	_, err = pidf.WriteString(string(os.Getpid()))
	if err != nil {
		zap.S().Fatalf("Failed to write pid: %s", err)
	}

	defer func() {
		pidf.Close()
		os.Remove("./.ninedns.pid")
	}()

	zap.S().Infof("Nine DNS started on %s", utils.C.Addr)

	serv := &dns.Server{Addr: utils.C.Addr, Net: "udp"}
	if err := serv.ListenAndServe(); err != nil {
		zap.S().Fatal(err)
	}
}
