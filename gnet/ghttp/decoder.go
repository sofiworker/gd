package ghttp

import (
	"encoding/json"
	"io"
)

type Decoder interface {
	Decode(reader io.Reader, v interface{}) error
}

func NewJsonDecoder() Decoder {
	return &jsonDecoder{}
}

type jsonDecoder struct{}

func (d *jsonDecoder) Decode(reader io.Reader, v interface{}) error {
	return json.NewDecoder(reader).Decode(v)
}

type xmlDecoder struct{}

func (d *xmlDecoder) Decode(reader io.Reader, v interface{}) error {
	return nil
}

type yamlDecoder struct {
}

func (d *yamlDecoder) Decode(reader io.Reader, v interface{}) error {
	return nil
}
