package database

import (
	"github.com/webhkp/godft/internal/consts"
	"github.com/webhkp/godft/pkg/util"
)

type Collection struct {
	Limit  int
	Sort   map[string]int
	Fields []string
}

func NewCollection(config interface{}) (collection *Collection) {
	collection = NewDefaultCollection()

	if config == nil {
		return
	}

	parsedConfig := config.(map[interface{}]interface{})

	if _, ok := parsedConfig[consts.LimitKey]; ok {
		collection.Limit = parsedConfig[consts.LimitKey].(int)
	}

	if _, ok := parsedConfig[consts.SortKey]; ok {
		for key, val := range util.ProcessInterfaceMap(parsedConfig[consts.SortKey].(map[interface{}]interface{})) {
			if val.(int) >= 0 {
				collection.Sort[key] = 1
				continue
			}
			collection.Sort[key] = -1
		}
	}

	if _, ok := parsedConfig[consts.FieldKey]; ok {
		for _, val := range parsedConfig[consts.FieldKey].([]interface{}) {
			collection.Fields = append(collection.Fields, val.(string))
		}
	}

	return
}

func NewDefaultCollection() *Collection {
	return &Collection{
		Limit:  consts.NoLimit,
		Sort:   map[string]int{},
		Fields: []string{},
	}
}
