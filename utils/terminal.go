package utils

import (
	"fmt"
	"os"

	"github.com/hkbarton/arango-cli/state"
)

func outputText(data []string) {
	if data != nil && len(data) > 0 {
		for _, line := range data {
			fmt.Println(line)
		}
	}
}

func outputError(err error) {
	fmt.Printf("\033[1;31m%s\033[0m\n", err)
	if _, ok := err.(FatalError); ok {
		os.Exit(1)
	}
}

// Output print command result on terminal
func Output(data interface{}) {
	switch data.(type) {
	case []string:
		outputText(data.([]string))
	case error:
		outputError(data.(error))
	}
	fmt.Print(state.Prompt())
}
