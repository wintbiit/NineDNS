//go:build mysql

package provider

import (
	"github.com/wintbiit/ninedns/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlProvider struct {
	Provider
	*gorm.DB
	dsn string
}

func init() {
	constructors["mysql"] = newMysqlProvider
}

func newMysqlProvider(dsn string) (Provider, error) {
	provider := &MysqlProvider{
		dsn: dsn,
	}

	db, err := gorm.Open(mysql.Open(provider.dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	provider.DB = db

	return provider, nil
}

func (p *MysqlProvider) Provide(ruleset string) ([]model.Record, error) {
	tx := p.Begin()
	defer tx.Rollback()

	var records []model.Record

	if err := tx.Table(ruleset).Find(&records).Error; err != nil {
		return nil, err
	}

	tx.Commit()

	return records, nil
}

func (p *MysqlProvider) AutoMigrate(table string) error {
	return p.DB.Table(table).AutoMigrate(&model.Record{})
}
