package utils

import (
	"bytes"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"sigs.k8s.io/yaml"
)

func ToYaml(it any) ([]byte, error) {
	var (
		bts []byte
		err error
	)

	buffer := bytes.NewBuffer(bts)
	enc := jsontext.NewEncoder(buffer, json.FormatNilMapAsNull(true), json.FormatNilSliceAsNull(true))

	if err = json.MarshalEncode(enc, it); err != nil {
		return bts, err
	}

	bts = buffer.Bytes()

	if bts, err = yaml.JSONToYAML(bts); err != nil {
		return bts, err
	}

	return bts, nil
}
