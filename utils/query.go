package utils

import (
	driver "github.com/arangodb/go-driver"
	"github.com/hkbarton/arango-cli/state"
)

type DocumentResult struct {
	ID       string
	key      string
	Document map[string]interface{}
}

func QueryOne(db driver.Database, aql string, args map[string]interface{}) (*DocumentResult, error) {
	if db == nil {
		db = state.CurrentDB()
	}
	cursor, err := db.Query(nil, aql, args)
	if err != nil {
		return nil, err
	}
	var result DocumentResult
	meta, err := cursor.ReadDocument(nil, &result.Document)
	if err != nil {
		return nil, err
	}
	result.ID = string(meta.ID)
	result.key = meta.Key
	return &result, nil
}
