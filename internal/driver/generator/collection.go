package generator

import (
	"github.com/webhkp/godft/internal/consts"
)

type Collection struct {
	Limit  int
	Fields consts.GeneratorCollectionFieldType
}

func NewCollection(config interface{}) (collection *Collection) {
	collection = NewDefaultCollection()

	if config == nil {
		return
	}

	parsedConfig := config.(map[interface{}]interface{})

	if _, ok := parsedConfig[consts.LimitKey]; ok {
		collection.Limit = parsedConfig[consts.LimitKey].(int)
	} else {
		collection.Limit = consts.GeneratorDefaultRows
	}

	if _, ok := parsedConfig[consts.FieldKey]; ok {
		collection.Fields = make(consts.GeneratorCollectionFieldType)

		for key, val := range parsedConfig[consts.FieldKey].(map[interface{}]interface{}) {
			collection.Fields[key.(string)] = val.(string)
		}
	}

	return
}

func NewDefaultCollection() *Collection {
	return &Collection{
		Limit:  consts.NoLimit,
		Fields: make(consts.GeneratorCollectionFieldType),
	}
}
