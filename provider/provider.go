package provider

import (
	"fmt"

	"github.com/wintbiit/ninedns/model"
)

var constructors = make(map[string]func(string) (Provider, error))

type Provider interface {
	Provide(ruleset string) ([]model.Record, error)
	AutoMigrate(table string) error
}

func NewProvider(name, config string) (Provider, error) {
	constructor, ok := constructors[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", name)
	}

	return constructor(config)
}
