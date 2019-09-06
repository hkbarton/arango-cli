package terminal

import (
	"fmt"

	"github.com/hkbarton/arango-cli/state"
)

// Output print command result on terminal
func Output(data []string) {
	if data != nil && len(data) > 0 {
		for _, line := range data {
			fmt.Println(line)
		}
	}
	fmt.Print(state.Prompt())
}
