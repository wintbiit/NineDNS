package provider

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/wintbiit/ninedns/model"
)

// FileProvider reads records from a file.
// record format like:
// pdm                     IN      A       172.30.162.35
type FileProvider struct {
	Provider
	file string
}

func newFileProvider(file string) (*FileProvider, error) {
	provider := &FileProvider{
		file: file,
	}

	return provider, nil
}

func (p *FileProvider) Provide(string) ([]model.Record, error) {
	f, err := os.Open(p.file)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	reader := bufio.NewReader(f)
	var records []model.Record

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}

		record, err := readLine(string(line))
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

func readLine(line string) (model.Record, error) {
	var record model.Record

	reg := regexp.MustCompile(`\s+`)
	fields := reg.Split(line, -1)
	if len(fields) < 4 {
		return record, fmt.Errorf("invalid record: %s", line)
	}

	record.Host = fields[0]
	record.Type = model.RecordType(fields[1])
	record.Value = model.RecordValue(fields[3])

	return record, nil
}

func (p *FileProvider) AutoMigrate(_ string) error {
	return nil
}
