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
  "addr": ":53",                            // listen address
  "debug": true,                            // debug mode
  "domains": [                              // dns server domains to resolve
    {
      "domain": "example.com",              // domain name
      "authoritative": true,                // authoritative mode
      "recursion": false,                   // recursion mode
      "upstream": "223.5.5.5:53",           // upstream dns server, available when recursion is true
      "ttl": 600,                           // default ttl
      "mysql": "root:luoqiwen20040602@tcp(localhost:3306)/dns", // dns record provider mysql dsn
      "rules": [                            // dns resolve rule set
        {
          "cidr": "127.0.0.1/24",           // remote addr cidr rule. Hits when remote addr matches
          "name": "example_com"             // rule set name, also table name in mysql
        }
      ]
    }
  ],
  "redis": {                                // redis config
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