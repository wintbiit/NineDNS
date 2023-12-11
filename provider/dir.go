package provider

import (
	"github.com/wintbiit/ninedns/model"
)

type DirProvider struct {
	Provider
	dir string
}

func init() {
	constructors["dir"] = newDirProvider
}

func newDirProvider(dir string) (Provider, error) {
	provider := &DirProvider{
		dir: dir,
	}

	return provider, nil
}

func (p *DirProvider) Provide(ruleset string) ([]model.Record, error) {
	provider, err := newFileProvider(p.dir + "/" + ruleset)
	if err != nil {
		return nil, err
	}

	return provider.Provide(ruleset)
}

func (p *DirProvider) AutoMigrate(_ string) error {
	return nil
}
