package model

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

type (
	RecordType  string
	RecordValue string
)

const (
	RecordTypeA     RecordType = "A"
	RecordTypeAAAA  RecordType = "AAAA"
	RecordTypeCNAME RecordType = "CNAME"
	RecordTypeMX    RecordType = "MX"
	RecordTypeNS    RecordType = "NS"
	RecordTypeSOA   RecordType = "SOA"
	RecordTypeSRV   RecordType = "SRV"
	RecordTypeTXT   RecordType = "TXT"
)

type Record struct {
	Host     string      `json:"host"`
	Type     RecordType  `json:"type"`
	Value    RecordValue `json:"value"`
	Weight   uint16      `json:"weight,default=1"`
	Disabled bool        `json:"disabled,default=false"`
	Note     string      `json:"-"`
}

type SOARecord struct {
	NS      string
	MBox    string
	Serial  uint32
	Refresh uint32
	Retry   uint32
	Expire  uint32
	MinTTL  uint32
}

type SRVRecord struct {
	Priority uint16
	Weight   uint16
	Port     uint16
	Target   string
}

type MXRecord struct {
	Preference uint16 `json:"preference"`
	MX         string `json:"mx"`
}

func (v *RecordValue) IP() net.IP {
	return net.ParseIP(string(*v))
}

func (v *RecordValue) String() string {
	return string(*v)
}

func (v *RecordValue) SOA() (*SOARecord, error) {
	sp := strings.Split(string(*v), " ")
	if len(sp) != 7 {
		return nil, fmt.Errorf("invalid SOA record: %s", string(*v))
	}

	var soa SOARecord
	soa.NS = sp[0]
	soa.MBox = sp[1]
	serial, err := strconv.ParseUint(sp[2], 10, 32)
	if err != nil {
		return nil, err
	}

	soa.Serial = uint32(serial)

	refresh, err := strconv.ParseUint(sp[3], 10, 32)
	if err != nil {
		return nil, err
	}

	soa.Refresh = uint32(refresh)

	retry, err := strconv.ParseUint(sp[4], 10, 32)
	if err != nil {
		return nil, err
	}

	soa.Retry = uint32(retry)

	expire, err := strconv.ParseUint(sp[5], 10, 32)
	if err != nil {
		return nil, err
	}

	soa.Expire = uint32(expire)

	minTTL, err := strconv.ParseUint(sp[6], 10, 32)
	if err != nil {
		return nil, err
	}

	soa.MinTTL = uint32(minTTL)

	return &soa, nil
}

func (v *RecordValue) SRV() (*SRVRecord, error) {
	sp := strings.Split(string(*v), " ")
	if len(sp) != 4 {
		return nil, fmt.Errorf("invalid SRV record: %s", string(*v))
	}

	var srv SRVRecord
	priority, err := strconv.ParseUint(sp[0], 10, 16)
	if err != nil {
		return nil, err
	}

	srv.Priority = uint16(priority)

	weight, err := strconv.ParseUint(sp[1], 10, 16)
	if err != nil {
		return nil, err
	}

	srv.Weight = uint16(weight)

	port, err := strconv.ParseUint(sp[2], 10, 16)
	if err != nil {
		return nil, err
	}

	srv.Port = uint16(port)
	srv.Target = sp[3]

	return &srv, nil
}

func (v *RecordValue) MX() (*MXRecord, error) {
	sp := strings.Split(string(*v), " ")
	if len(sp) != 2 {
		return nil, fmt.Errorf("invalid MX record: %s", string(*v))
	}

	var mx MXRecord
	preference, err := strconv.ParseUint(sp[0], 10, 16)
	if err != nil {
		return nil, err
	}

	mx.Preference = uint16(preference)
	mx.MX = sp[1]

	return &mx, nil
}

func ReadRecordType(t uint16) RecordType {
	return RecordType(dns.TypeToString[t])
}

func (t RecordType) DnsType() uint16 {
	switch t {
	case RecordTypeA:
		return dns.TypeA
	case RecordTypeAAAA:
		return dns.TypeAAAA
	case RecordTypeCNAME:
		return dns.TypeCNAME
	case RecordTypeMX:
		return dns.TypeMX
	case RecordTypeNS:
		return dns.TypeNS
	case RecordTypeSOA:
		return dns.TypeSOA
	case RecordTypeSRV:
		return dns.TypeSRV
	case RecordTypeTXT:
		return dns.TypeTXT
	default:
		return dns.TypeNone
	}
}

func (t RecordType) String() string {
	return string(t)
}
