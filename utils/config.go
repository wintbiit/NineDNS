package utils

import (
	"flag"
	"os"

	"github.com/wintbiit/ninedns/model"
)

var C *model.Config

func init() {
	config := flag.String("config", "./ninedns.json", "config file path")
	flag.Parse()

	f, err := os.Open(*config)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var conf model.Config
	if err = UnmarshalFromReader(f, &conf); err != nil {
		panic(err)
	}

	C = &conf
}
