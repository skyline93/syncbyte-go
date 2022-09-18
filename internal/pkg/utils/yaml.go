package utils

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

func DecodeFromFile(file string, v interface{}) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, v)
	if err != nil {
		return err
	}

	return nil
}

func EncodeToMap(in interface{}) (map[string]interface{}, error) {
	out, err := yaml.Marshal(in)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	if err = yaml.Unmarshal(out, &m); err != nil {
		return nil, err
	}

	return m, nil
}

func EncodeToBytes(in interface{}) ([]byte, error) {
	out, err := yaml.Marshal(in)
	if err != nil {
		return nil, err
	}

	return out, err
}
