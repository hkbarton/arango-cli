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

// Commander run command and return result
type Commander interface {
	Run()
}

func (c Command) String() string {
	return fmt.Sprintf("action: %s, args: %v, options: %v", c.action, c.args, c.options)
}

// ParseCommand parse the input of string slice to Command
func ParseCommand(input []string) (*Command, error) {
	if input == nil || len(input) < 1 {
		return nil, fmt.Errorf("no input for parse")
	}
	command := &Command{
		action:  input[0],
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
