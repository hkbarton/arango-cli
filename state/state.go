package state

import "fmt"

var state = map[string]string{
	"currentHost": "",
	"currentDB":   "",
}

// SetState set global state
func SetState(newState map[string]string) {
	for key, value := range newState {
		state[key] = value
	}
}

// Prompt render terminal prompt string by state
func Prompt() string {
	return fmt.Sprintf("\033[1;34m%s\033[0m.\033[32m%s\033[0m > ",
		state["currentHost"], state["currentDB"])
}
