//go:build !nolark

package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/wintbiit/larki"
	"github.com/wintbiit/ninedns/model"
)

type LarkProvider struct {
	Provider
	*larki.Client
	baseId string
}

func init() {
	constructors["lark"] = newLarkProvider
}

func newLarkProvider(config string) (Provider, error) {
	sp := strings.Split(config, " ")
	if len(sp) != 3 {
		return nil, fmt.Errorf("invalid lark config: %s you should use `${appId} ${appSecret} ${baseId}` format", config)
	}

	client, err := larki.NewClient(sp[0], sp[1], "", "")
	if err != nil {
		return nil, err
	}

	return &LarkProvider{
		Client: client,
		baseId: sp[2],
	}, nil
}

func (p *LarkProvider) Provide(ruleset string) ([]model.Record, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	tables, err := p.ListBaseTables(ctx, p.baseId)
	if err != nil {
		return nil, err
	}

	var tableId string
	for _, table := range tables {
		if *table.Name == ruleset {
			tableId = *table.TableId
			break
		}
	}

	if tableId == "" {
		return nil, fmt.Errorf("table %s not found", ruleset)
	}

	baseRec, err := p.GetRecords(ctx, p.baseId, tableId, "", 0)
	if err != nil {
		return nil, err
	}

	var records []model.Record
	for _, rec := range baseRec {
		weight, err := strconv.ParseUint(rec.Fields["Weight"].(string), 10, 16)
		if err != nil {
			return nil, err
		}

		disabled, ok := rec.Fields["Disabled"].(bool)
		if !ok {
			disabled = false
		}

		records = append(records, model.Record{
			Host:     rec.Fields["Host"].(string),
			Type:     model.RecordType(rec.Fields["Type"].(string)),
			Value:    model.RecordValue(rec.Fields["Value"].(string)),
			Weight:   uint16(weight),
			Disabled: disabled,
		})
	}

	return records, nil
}
