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
Mysql Provider connects to a mysql database and reads records from it. Table name is rule set name.
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

### SQLite Provider
config:
```json
{
  "provider": {
    "sqlite": "sqlite3://path/to/sqlite.db"
  }
}
```
Like Mysql Provider, SQLite Provider connects to a sqlite database and reads records from it. Table name is rule set name.
SQL is re-read according to ttl

> Using `github.com/glebarez/sqlite` as sqlite driver, please refer to this repo for more usage.

### File Provider
config:
```json
{
  "provider": {
    "file": "path/to/file.db"
  }
}
```
File Provider reads records from a file. File format is standard zone file format.
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
Dir Provider reads records from a directory. File format is standard zone file format.
ruleset name is used as file name.
files in the directory will be re-read according to ttl