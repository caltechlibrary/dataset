package dataset

import (
	"bytes"

	// 3rd Party Package
	"gopkg.in/yaml.v3"
)

// YAMLUnmarshal is a custom YAML decoder so we can treat numbers easier
func YAMLUnmarshal(src []byte, data interface{}) error {
	r := bytes.NewBuffer(src)
	dec := yaml.NewDecoder(r)
	dec.KnownFields(true)
	if err := dec.Decode(data); err != nil {
		return err
	}
	return nil
}

// YAMLMarshal provides provide a custom json encoder to solve a an issue with
// HTML entities getting converted to UTF-8 code points by json.Marshal(), json.MarshalIndent().
func YAMLMarshal(data interface{}) ([]byte, error) {
	buf := []byte{}
	w := bytes.NewBuffer(buf)
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	src := w.Bytes()
	return src, err
}

// YAMLMarshalIndent provides provide a custom json encoder to solve a an issue with
// HTML entities getting converted to UTF-8 code points by json.Marshal(), json.MarshalIndent().
func YAMLMarshalIndent(data interface{}, spaces int) ([]byte, error) {
	buf := []byte{}
	w := bytes.NewBuffer(buf)
	enc := yaml.NewEncoder(w)
	enc.SetIndent(spaces)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	src := w.Bytes()
	return src, err
}
