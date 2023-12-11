package model

import (
	"github.com/miekg/dns"
)

type RuleProvider interface {
	FindRecord(name string, t uint16) *Record
	FindRecords(name string, t uint16) []Record
	Header(r *Record) dns.RR_Header
	Recursion() bool
	Exchange(r *dns.Msg) (*dns.Msg, error)
}
