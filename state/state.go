package state

import (
	"fmt"
	"sync"

	driver "github.com/arangodb/go-driver"
)

var state = map[string]interface{}{
	"currentHost": "",
	"currentDB":   "",
	"dbClient":    nil,
}

var mu sync.Mutex

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

// Prompt render terminal prompt string by state
func Prompt() string {
	return fmt.Sprintf("\033[1;34m%s\033[0m.\033[32m%s\033[0m > ",
		state["currentHost"], state["currentDB"])
}
