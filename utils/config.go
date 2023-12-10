package utils

import (
	"flag"
	"os"

	"github.com/wintbiit/ninedns/model"
	"go.uber.org/zap"
)

var C *model.Config

func init() {
	config := flag.String("config", "./ninedns.json", "config file path")
	flag.Parse()

	f, err := os.Open(*config)
	if err != nil {
		zap.S().Fatalf("Failed to open config file: %s", err)
	}

	var conf model.Config
	if err = UnmarshalFromReader(f, &conf); err != nil {
		zap.S().Fatalf("Failed to parse config file: %s", err)
	}

	C = &conf
	zap.S().Debugf("Loaded config: %+v", C)
}
