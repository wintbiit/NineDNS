//go:build lark

package utils

import (
	"io"

	"github.com/bytedance/sonic"
)

type JsonSonic struct {
	sonic.API
}

func (j JsonSonic) UnmarshalFromReader(reader io.Reader, v interface{}) error {
	return j.NewDecoder(reader).Decode(v)
}

func init() {
	json = JsonSonic{
		API: sonic.ConfigFastest,
	}
}
