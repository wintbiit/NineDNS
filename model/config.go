package model

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
	Rules         map[string]Rule `json:"rules"`
	Authoritative bool            `json:"authoritative,default=true"`
	Recursion     bool            `json:"recursion,default=false"`
	Upstream      string          `json:"upstream,default=127.0.0.1:53"`
	MySQL         string          `json:"mysql"`
	TTL           uint32          `json:"ttl,default=60"`
	Domain        string          `json:"-"`
}

type Rule struct {
	CIDRs   []string `json:"cidrs"`
	Name    string   `json:"-"`
	Records []Record `json:"records"`
}
