package provider

import (
	"fmt"

	"github.com/wintbiit/ninedns/model"
)

type Provider interface {
	Provide(ruleset string) ([]model.Record, error)
	AutoMigrate(table string) error
}

func NewProvider(name, config string) (Provider, error) {
	switch name {
	case "mysql":
		return newMysqlProvider(config)
	case "sqlite":
		return newSQLiteProvider(config)
	case "file":
		return newFileProvider(config)
	case "dir":
		return newDirProvider(config)
	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
}
