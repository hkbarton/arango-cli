package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

const maxCommandHistory = 200

var (
	commandHistory []string
	historyCursor  int
)

func getCommandHistoryFile() string {
	var folder string
	curUser, err := user.Current()
	if err == nil {
		folder = curUser.HomeDir
	} else {
		folder, _ = filepath.Abs(os.Args[0])
	}
	historyFile := filepath.Join(folder, ".arango-cli-his")
	return historyFile
}

func init() {
	historyFile := getCommandHistoryFile()
	content, err := ioutil.ReadFile(getCommandHistoryFile())
	if err != nil {
		_, err = os.Create(historyFile)
		if err != nil {
			panic(err)
		}
		commandHistory = make([]string, 0, 20)
	} else {
		commandHistory = strings.Split(string(content), "\n")
	}
	historyCursor = -1
}

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

// PushHistory pushes command input into history list
func PushHistory(commandStr string) {
	commandHistory = append(commandHistory, strings.TrimSpace(commandStr))
	if len(commandHistory) > maxCommandHistory {
		start := len(commandHistory) - maxCommandHistory
		commandHistory = commandHistory[start:]
	}
	data := []byte(strings.TrimSpace(strings.Join(commandHistory, "\n")))
	ioutil.WriteFile(getCommandHistoryFile(), data, 0755)
}

// Last return last command from history list
func Last() string {
	if len(commandHistory) == 0 {
		return ""
	}
	historyCursor--
	if historyCursor < 0 {
		historyCursor = len(commandHistory) - 1
	}
	return commandHistory[historyCursor]
}

// Next return next command from history list
func Next() string {
	if len(commandHistory) == 0 {
		return ""
	}
	historyCursor++
	if historyCursor > len(commandHistory)-1 {
		historyCursor = 0
	}
	return commandHistory[historyCursor]
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
	timeout, _ := time.ParseDuration("30s")
	go runner.Run(command, resultChan)
	select {
	case result := <-resultChan:
		return result
	case <-time.After(timeout):
		return errors.New("Connection timeout")
	}
}

// RunByAction automatically pick command by action and run it
func RunByAction(command *Command) interface{} {
	runner, exists := commandActionMap[command.Action]
	if !exists {
		return fmt.Errorf("Unknown command: %s", command.Action)
	}
	return Run(command, runner)
}
