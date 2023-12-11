//go:build postgres

package provider

import (
	"github.com/wintbiit/ninedns/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresProvider struct {
	Provider
	*gorm.DB
	dsn string
}

func init() {
	constructors["postgres"] = newPostgresProvider
}

func newPostgresProvider(dsn string) (Provider, error) {
	provider := &PostgresProvider{
		dsn: dsn,
	}

	db, err := gorm.Open(postgres.Open(provider.dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	provider.DB = db

	return provider, nil
}

func (p *PostgresProvider) Provide(ruleset string) ([]model.Record, error) {
	tx := p.Begin()
	defer tx.Rollback()

	var records []model.Record

	if err := tx.Table(ruleset).Find(&records).Error; err != nil {
		return nil, err
	}

	tx.Commit()

	return records, nil
}

func (p *PostgresProvider) AutoMigrate(table string) error {
	return p.DB.Table(table).AutoMigrate(&model.Record{})
}
