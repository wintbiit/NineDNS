## NineDNS Record Provider

### Mysql Provider
config:
```json
{
  "provider": {
    "mysql": "mysql://user:password@host:port/dbname?charset=utf8&parseTime=True&loc=Local"
  }
}
```
`Mysql Provider` connects to a mysql database and reads records from it. Table name is rule set name.
SQL is re-read according to ttl

Table schema:
```sql
CREATE TABLE if not exists `rule_set_name` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `host` varchar(255) NOT NULL DEFAULT '',
  `type` varchar(255) NOT NULL DEFAULT '',
  `value` varchar(255) NOT NULL DEFAULT '',
  `ttl` int(11) unsigned NOT NULL DEFAULT '0',
  `weight` int(11) unsigned NOT NULL DEFAULT '0',
  `disabled` tinyint(1) unsigned NOT NULL DEFAULT '0',
  `note` varchar(255) NOT NULL DEFAULT '',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`host`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

Not include in default build, you need to add `mysql` tag to build it.
```bash
go build -tags "mysql"
```

### SQLite Provider
config:
```json
{
  "provider": {
    "sqlite": "sqlite3://path/to/sqlite.db"
  }
}
```
Like `Mysql Provider`, `SQLite Provider` connects to a sqlite database and reads records from it. Table name is rule set name.
SQL is re-read according to ttl

> Using `github.com/glebarez/sqlite` as sqlite driver, please refer to this repo for more usage.

Disabled by default, you need to add `sqlite` tag to build it.
```bash
go build -tags "sqlite"
```

### File Provider
config:
```json
{
  "provider": {
    "file": "path/to/file.db"
  }
}
```
`File Provider` reads records from a file. File format is standard zone file format.
File Provider does not support ruleset variant, if you want to use variant, please use [Dir Provider](#dir-provider).
File will be re-read according to ttl

### Dir Provider
config:
```json
{
  "provider": {
    "dir": "path/to/dir"
  }
}
```
`Dir Provider` reads records from a directory. File format is standard zone file format.
ruleset name is used as file name.
files in the directory will be re-read according to ttl



### Lark Provider
config:
```json
{
  "provider": {
    "lark": "cli_xxx xxx xxx"
  }
}
```
`Lark Provider` reads records from lark bitable. Also, ruleset name is used as table name.

#### Minimize binary size
 Please note, lark provider introduced [oapi-lark-go](https://github.com/larksuite/oapi-sdk-go) and [sonic](https://github.com/bytedance/sonic),
 which largely increases binary size, it's disabled by default. You need to add `lark` tag to build it.
 ```bash
go build -tags "lark"
```

### Postgres Provider
config:
```json
{
  "provider": {
    "postgres": "postgres://user:password@host:port/dbname?sslmode=disable"
  }
}
```

`Postgres Provider` connects to a postgres database and reads records from it. Table name is rule set name.

Disabled by default, you need to add `postgres` tag to build it.
```bash
go build -tags "postgres"
```

