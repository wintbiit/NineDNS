# Nine DNS
[![Build](https://github.com/wintbiit/NineDNS/actions/workflows/build.yml/badge.svg)](https://github.com/wintbiit/NineDNS/actions/workflows/build.yml)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/wintbiit/NineDNS)](https://github.com/wintbiit/NineDNS/releases)

Nine DNS is a flexible and high-performance dns server.

## Why NineDNS
### 1. ACL DNS Resolve
`NineDNS` aims to differentiated dns resolve based on question source.

### 2. DNS Records from Remote
Meanwhile `NineDNS` supports providing dns records from remote databases such as `MySQL` or `PostgreSQL`

### 3. Cloud Native
`NineDNS` can be part of cloud native component. `NineDNS` supports cache sharing, load balancing and log tracing.

## Usage
Define a config:
```json
{
  "addr": ":53",                   // listen address
  "debug": true,                   // debug mode
  "domains": {                     // dns resolve domain key-value pairs. domain <===> resolve config
    "example.com": {
      "authoritative": true,       // authoritative mode
      "recursion": false,          // recursion mode
      "upstream": "223.5.5.5:53",  // upstream dns server, only works in recursion mode
      "ttl": 600,                  // default ttl, attention: ttl is server level, not record level. server re-fetch record source ttl
      "mysql": "root:123456@tcp(localhost:3306)/dns", // mysql dns record source
      "rules": {                   // dns resolve match rules. name <===> rule. Name is also used as table name in mysql record source
        "all": {
          "cidrs": [               // cidr match
            "0.0.0.0/0"
          ]
        }
      },
      "records": [                 // manually defined records. These records would overrider others
        {
          "host": "git",        
          "type": "CNAME",
          "value": "qcloud.example.com"
        }
      ]
    }
  },
  "redis": {                      // redis config
    "addr": "localhost:6379",
    "db": 8
  }
}
```

And that's all. Run `NineDNS` with config file:
```shell
$ ninedns -c config.json
```
> `NineDNS` autoloads `ninedns.json` in current directory if `-c` is not specified.