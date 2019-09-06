package commands

import (
	"fmt"
	"strings"
)

// Command represent command to be execute
type Command struct {
	action  string
	args    []string
	options map[string]string
}

func (c Command) String() string {
	return fmt.Sprintf("action: %s, args: %v, options: %v", c.action, c.args, c.options)
}

// CommandRunner run command and hanele output
type CommandRunner interface {
	Run(c *Command, resultChan chan []string)
}

var commandActionMap map[string]CommandRunner

func init() {
	commandActionMap = map[string]CommandRunner{
		"exit": new(ExitCommandRunner),
	}
}

// ParseCommand parse the input of string slice to Command
func ParseCommand(input []string) (*Command, error) {
	if input == nil || len(input) < 1 {
		return nil, fmt.Errorf("no input for parse")
	}
	command := &Command{
		action:  strings.TrimSpace(input[0]),
		args:    make([]string, 0, 2),
		options: make(map[string]string),
	}
	for i := 1; i < len(input); i++ {
		seg := input[i]
		if strings.HasPrefix(seg, "-") {
			optionKey := strings.TrimPrefix(seg, "-")
			// input starts with - is options
			if i > len(input)-2 {
				// last input
				command.options[optionKey] = ""
			} else {
				command.options[optionKey] = input[i+1]
				i++
			}
		} else {
			// args
			command.args = append(command.args, seg)
		}
	}
	return command, nil
}

// RunCommand run command and hanele result
func RunCommand(command *Command, runner CommandRunner) []string {
	resultChan := make(chan []string)
	go runner.Run(command, resultChan)
	result := <-resultChan
	return result
}

// RunCommandByAction automatically pick command by action and run it
func RunCommandByAction(command *Command) []string {
	runner, exists := commandActionMap[command.action]
	if !exists {
		return []string{fmt.Sprintf("Unknown command: %s", command.action)}
	}
	return RunCommand(command, runner)
}
