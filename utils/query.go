package utils

import (
	"context"
	"encoding/json"
	"fmt"

	driver "github.com/arangodb/go-driver"
	"github.com/hkbarton/arango-cli/state"
)

type DocumentResult struct {
	ID       string
	key      string
	Document map[string]interface{}
}

func (d DocumentResult) String() string {
	docData, err := json.Marshal(d.Document)
	if err != nil {
		return fmt.Sprintf("%s\t\t%v", d.ID, "Content Parsing Error")
	}
	return fmt.Sprintf("%s\t%v", d.ID, string(docData))
}

func (d DocumentResult) JsonString() string {
	docData, err := json.MarshalIndent(d.Document, "", "\t")
	if err != nil {
		return fmt.Sprintf("%s\t\t%v", d.ID, "Content Parsing Error")
	}
	return string(docData)
}

func getCursor(db driver.Database, aql string, args map[string]interface{}) (driver.Cursor, error) {
	if db == nil {
		db = state.CurrentDB()
	}
	ctx := driver.WithQueryCount(context.Background())
	cursor, err := db.Query(ctx, aql, args)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

func QueryOne(db driver.Database, aql string, args map[string]interface{}) (*DocumentResult, error) {
	cursor, err := getCursor(db, aql, args)
	defer cursor.Close()
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

func QueryAll(db driver.Database, aql string, args map[string]interface{}) ([]*DocumentResult, int, error) {
	cursor, err := getCursor(db, aql, args)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close()
	count := cursor.Count()
	results := make([]*DocumentResult, 0, 10)
	for cursor.HasMore() {
		var doc DocumentResult
		meta, err := cursor.ReadDocument(nil, &doc.Document)
		if err != nil {
			return nil, 0, err
		}
		doc.ID = string(meta.ID)
		doc.key = meta.Key
		results = append(results, &doc)
	}
	return results, int(count), nil
}
