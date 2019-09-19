package utils

import (
	"fmt"
	"os"

	"github.com/hkbarton/arango-cli/state"
)

const (
	RedColor = iota
	BlueColor
	GreenColor
	CyanColor
)

func outputText(data []string) {
	if data != nil && len(data) > 0 {
		fmt.Println()
		for _, line := range data {
			fmt.Println(line)
		}
		fmt.Println()
	}
}

func outputError(err error) {
	errStr := fmt.Sprintf("%s", err)
	fmt.Println(Color(errStr, RedColor))
	if _, ok := err.(FatalError); ok {
		os.Exit(1)
	}
}

// Color and ANSI color to the input string
func Color(input string, colorCode int) string {
	switch colorCode {
	case RedColor:
		return fmt.Sprintf("\033[1;31m%s\033[0m", input)
	case BlueColor:
		return fmt.Sprintf("\033[1;34m%s\033[0m", input)
	case GreenColor:
		return fmt.Sprintf("\033[1;32m%s\033[0m", input)
	case CyanColor:
		return fmt.Sprintf("\033[1;96m%s\033[0m", input)
	}
	return input
}

// Output print command result on terminal
func Output(data interface{}) {
	switch data.(type) {
	case []string:
		outputText(data.([]string))
	case error:
		outputError(data.(error))
	}
	fmt.Print(CurrentPrompt())
}

// CurrentPrompt render terminal prompt string by state
func CurrentPrompt() string {
	return fmt.Sprintf("%s.%s > ",
		Color(state.GetState("currentHost").(string), BlueColor),
		Color(state.GetState("currentDB").(string), GreenColor))
}
