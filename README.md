# Nine DNS
[![Build](https://github.com/wintbiit/NineDNS/actions/workflows/build.yml/badge.svg)](https://github.com/wintbiit/NineDNS/actions/workflows/build.yml)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/wintbiit/NineDNS)](https://github.com/wintbiit/NineDNS/releases)

Nine DNS is a flexible DNS server that offers DNS resolution based on the question source.

## Why NineDNS
### 1. ACL DNS Resolve
`NineDNS` aims to provide a flexible way to resolve DNS records.

You can match different question source by cidr, port, protocol, and so on.

Clients can use different dns resolve based on their network environment.

For example, you can filter clients by cidr, and resolve different dns records for them.

### 2. DNS Records from Remote
Moreover, NineDNS supports retrieving DNS records from remote databases such as MySQL or PostgreSQL.

It's easy to manage DNS records in a centralized way.

### 3. Cloud Native
`NineDNS` can integrate as part of cloud-native components. It supports cache sharing, load balancing, and log tracing.

## Main Features
### ACL DNS Resolve
The reason why we create `NineDNS` is to provide a flexible way to resolve DNS records according to the question source.

For example, you can match different question source by cidr, port, protocol, and so on.

Clients can use different dns resolve based on their network environment.

#### In production uses, for example:
You have your server deployed intranet and exposed server port via a jump server or tunnel.

And you want to resolve server domain to intranet ip when clients are in intranet, and resolve to public ip when clients are outside.

This is a typical use case of `NineDNS`. You simply `NS` your domain to `NineDNS` server, and configure `NineDNS` to resolve domain based on client's network environment, cidr for example.
![whiteboard_exported_image (1).png](images%2Fwhiteboard_exported_image%20%281%29.png)
### DNS Records from Remote
To be a full-featured DNS server, `NineDNS` supports retrieving DNS records from remote databases such as MySQL, files, lark, and so on.

Just like any Cloud DNS Providers, you can manage DNS records easily in a centralized way.

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
      "providers": {               // record source providers. Read [Record Provider](#record-provider) for more details
        "mysql": "root:123456@tcp(localhost:3306)/dns",
        "sqlite": "dns.db"
      },
      "rules": {                   // dns resolve match rules. name <===> rule. Name is also used as table name in mysql record source
        "all": {
          "cidrs": [               // cidr match
            "0.0.0.0/0"
          ],
          "ports": [
            "1-65535"              // port match
          ]
        }
      }
    }
  },
  "redis": {                      // redis config
    "addr": "localhost:6379",
    "db": 8
  }
}
```
> When a dns question comes in, `NineDNS` will find a rule defined that both matches the domain name and the rules.
> Then look for records in the record source defined in the rule.

Read [Record Provider](provider/README.md) for more details about `providers`.

And that's all! Run `NineDNS` with config file now:
```shell
$ ninedns -c config.json
```
> `NineDNS` auto-loads `ninedns.json` in current directory if `-c` is not specified.


## Downloads
Download from [releases](https://github.com/wintbiit/NineDNS/releases) page.

| Name           | Description                                           |
|----------------|-------------------------------------------------------|
| `ninedns-mini` | NineDNS binary without most providers.                |
| `ninedns`      | NineDNS binary with mysql provider and file providers |
| `ninedns-full` | NineDNS binary with all providers.                    |

## Special Thanks
Special thanks to miekg's graceful and well-designed go dns library [miekg/dns](https://github.com/miekg/dns)