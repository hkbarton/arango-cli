package state

import (
	"errors"
	"sync"

	driver "github.com/arangodb/go-driver"
)

var state = map[string]interface{}{
	"currentHost": "",
	"currentDB":   "",
	"dbClient":    nil,
}

var mu sync.Mutex
var currentDB driver.Database

// SetState set global state
func SetState(newState map[string]interface{}) {
	mu.Lock()
	defer mu.Unlock()
	for key, value := range newState {
		state[key] = value
	}
}

// GetState get state by key
func GetState(key string) interface{} {
	mu.Lock()
	defer mu.Unlock()
	return state[key]
}

// DBClient return db client
func DBClient() driver.Client {
	client := GetState("dbClient")
	return client.(driver.Client)
}

// SetCurrentDB sets current database
func SetCurrentDB(dbName string) error {
	mu.Lock()
	defer mu.Unlock()
	if dbName == "" {
		return errors.New("Must specifiy databse to be set")
	}
	if state["currentDB"].(string) == dbName {
		return nil
	}
	client := state["dbClient"].(driver.Client)
	db, err := client.Database(nil, dbName)
	if err == nil {
		state["currentDB"] = dbName
		currentDB = db
	}
	return err
}

// CurrentDB returns current database object
func CurrentDB() driver.Database {
	return currentDB
}
