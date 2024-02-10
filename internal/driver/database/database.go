package database

import (
	"github.com/webhkp/godft/internal/consts"
)

type Database struct {
	Connection    string
	ReadOnly      bool
	DatabaseName  string
	Collections   map[string]Collection
	AllCollection bool
}

func NewDatabase(config consts.FieldsType) (database *Database) {
	database = &Database{
		ReadOnly:      true,
		AllCollection: false,
		Collections:   map[string]Collection{},
	}

	if _, ok := config[consts.ConnectionKey]; ok {
		database.Connection = config[consts.ConnectionKey].(string)
	}

	if _, ok := config[consts.ReadOnlyKey]; ok {
		database.ReadOnly = config[consts.ReadOnlyKey].(bool)
	}

	if _, ok := config[consts.DatabaseKey]; ok {
		database.DatabaseName = config[consts.DatabaseKey].(string)
	} else {
		panic("Database name is required")
	}

	if _, ok := config[consts.CollectionKey]; ok {
		if config[consts.CollectionKey] == consts.AllCollectionValue {
			database.AllCollection = true
		} else if config[consts.CollectionKey] != nil {
			database.Collections = make(map[string]Collection)

			for key, collection := range config[consts.CollectionKey].(map[interface{}]interface{}) {
				database.Collections[key.(string)] = *NewCollection(collection)
			}
		}
	}

	return
}
