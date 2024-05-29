package utils

import (
	"bytes"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	ccyaml "github.com/goccy/go-yaml"
)

func ToYaml(it any) ([]byte, error) {
	var (
		bts []byte
		err error
	)

	buffer := bytes.NewBuffer(bts)
	enc := jsontext.NewEncoder(buffer, 
		json.FormatNilMapAsNull(true), 
		json.FormatNilSliceAsNull(true),
	)

	if err = json.MarshalEncode(enc, it); err != nil {
		return bts, err
	}

	bts = buffer.Bytes()

	if bts, err = ccyaml.JSONToYAML(bts); err != nil {
		return bts, err
	}

	return bts, nil
}
