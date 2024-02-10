package config

import (
	"os"

	"github.com/webhkp/godft/internal/consts"
	"gopkg.in/yaml.v2"
)

func Read(path string) consts.FlowData {
	data, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	// create a person struct and deserialize the data into that struct
	var fieldMap map[string]interface{}

	if err := yaml.Unmarshal(data, &fieldMap); err != nil {
		panic(err)
	}

	return fieldMap
}
