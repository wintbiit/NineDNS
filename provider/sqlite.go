package provider

import (
	"github.com/glebarez/sqlite"
	"github.com/wintbiit/ninedns/model"
	"gorm.io/gorm"
)

type SQLiteProvider struct {
	Provider
	*gorm.DB
	file string
}

func init() {
	constructors["sqlite"] = newSQLiteProvider
}

func newSQLiteProvider(file string) (Provider, error) {
	provider := &SQLiteProvider{
		file: file,
	}

	db, err := gorm.Open(sqlite.Open(provider.file), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	provider.DB = db

	return provider, nil
}

func (p *SQLiteProvider) Provide(string) ([]model.Record, error) {
	tx := p.Begin()
	defer tx.Rollback()

	var records []model.Record

	if err := tx.Find(&records).Error; err != nil {
		return nil, err
	}

	tx.Commit()

	return records, nil
}

func (p *SQLiteProvider) AutoMigrate(table string) error {
	return p.DB.Table(table).AutoMigrate(&model.Record{})
}
