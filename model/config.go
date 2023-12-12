package model

import (
	"strconv"
	"strings"
)

type Config struct {
	Addr    string            `json:"addr"`
	Debug   bool              `json:"debug"`
	Domains map[string]Domain `json:"domains"`
	Redis   struct {
		Addr     string `json:"addr"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	} `json:"redis"`
}

type Domain struct {
	Rules         map[string]Rule   `json:"rules"`
	Authoritative bool              `json:"authoritative,default=true"`
	Recursion     bool              `json:"recursion,default=false"`
	Upstream      string            `json:"upstream,default=127.0.0.1:53"`
	Providers     map[string]string `json:"providers"`
	TTL           uint32            `json:"ttl,default=60"`
	Tsig          *TSIG             `json:"tsig"`
}

type TSIG struct {
	Key string `json:"key"`
}

type Rule struct {
	CIDRs []string   `json:"cidrs"`
	Ports []PortRule `json:"ports"`
	Types []string   `json:"types"`
}

type PortRule string

func (p *PortRule) Contains(port int) bool {
	if strings.Count(string(*p), "-") == 1 {
		ports := strings.Split(string(*p), "-")
		if len(ports) != 2 {
			return false
		}

		start, err := strconv.Atoi(ports[0])
		if err != nil {
			return false
		}

		end, err := strconv.Atoi(ports[1])
		if err != nil {
			return false
		}

		return port >= start && port <= end
	}

	por, err := strconv.Atoi(string(*p))
	if err == nil {
		return por == port
	}

	return false
}
