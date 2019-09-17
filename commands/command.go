package commands

import (
	"fmt"
	"strings"
)

// Command represent command to be execute
type Command struct {
	Action  string
	Args    []string
	Options map[string]string
}

func (c Command) String() string {
	return fmt.Sprintf("action: %s, args: %v, options: %v", c.Action, c.Args, c.Options)
}

// Runner run command and hanele output
type Runner interface {
	Run(c *Command, resultChan chan interface{})
}

var commandActionMap map[string]Runner

func init() {
	commandActionMap = map[string]Runner{
		"exit": new(ExitCommandRunner),
		"ls":   new(ListCommandRunner),
		"use":  new(UseCommandRunner),
		"open": new(OpenCommandRunner),
		"info": new(InfoCommandRunner),
	}
}

// Parse parse the input of string slice to Command
func Parse(input []string) (*Command, error) {
	if input == nil || len(input) < 1 {
		return nil, fmt.Errorf("no input for parse")
	}
	command := &Command{
		Action:  strings.TrimSpace(input[0]),
		Args:    make([]string, 0, 2),
		Options: make(map[string]string),
	}
	for i := 1; i < len(input); i++ {
		seg := strings.TrimSpace(input[i])
		if strings.HasPrefix(seg, "-") {
			optionKey := strings.TrimPrefix(seg, "-")
			// input starts with - is options
			if i > len(input)-2 {
				// last input
				command.Options[optionKey] = ""
			} else {
				command.Options[optionKey] = input[i+1]
				i++
			}
		} else {
			// args
			command.Args = append(command.Args, seg)
		}
	}
	return command, nil
}

// Run run command and hanele result
func Run(command *Command, runner Runner) interface{} {
	resultChan := make(chan interface{})
	go runner.Run(command, resultChan)
	result := <-resultChan
	return result
}

// RunByAction automatically pick command by action and run it
func RunByAction(command *Command) interface{} {
	runner, exists := commandActionMap[command.Action]
	if !exists {
		return fmt.Errorf("Unknown command: %s", command.Action)
	}
	return Run(command, runner)
}
