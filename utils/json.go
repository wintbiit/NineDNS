package utils

import (
	"io"

	"github.com/bytedance/sonic"
)

var json = sonic.ConfigFastest

func MarshalToString(v interface{}) (string, error) {
	return json.MarshalToString(v)
}

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

func UnmarshalFromString(str string, v interface{}) error {
	return json.UnmarshalFromString(str, v)
}

func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func UnmarshalFromReader(reader io.Reader, v interface{}) error {
	return json.NewDecoder(reader).Decode(v)
}
