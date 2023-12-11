package utils

import (
	"io"

	jsonstd "encoding/json"
)

type API interface {
	MarshalToString(v interface{}) (string, error)
	Marshal(v interface{}) ([]byte, error)
	MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
	UnmarshalFromString(str string, v interface{}) error
	Unmarshal(data []byte, v interface{}) error
	UnmarshalFromReader(reader io.Reader, v interface{}) error
}

type JsonStd struct{}

var json API = JsonStd{}

func (JsonStd) MarshalToString(v interface{}) (string, error) {
	bytes, err := jsonstd.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (JsonStd) Marshal(v interface{}) ([]byte, error) {
	return jsonstd.Marshal(v)
}

func (JsonStd) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return jsonstd.MarshalIndent(v, prefix, indent)
}

func (JsonStd) UnmarshalFromString(str string, v interface{}) error {
	return jsonstd.Unmarshal([]byte(str), v)
}

func (JsonStd) Unmarshal(data []byte, v interface{}) error {
	return jsonstd.Unmarshal(data, v)
}

func (JsonStd) UnmarshalFromReader(reader io.Reader, v interface{}) error {
	return jsonstd.NewDecoder(reader).Decode(v)
}

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
	return json.UnmarshalFromReader(reader, v)
}
