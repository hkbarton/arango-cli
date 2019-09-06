package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hkbarton/arango-cli/commands"
	"github.com/hkbarton/arango-cli/state"
	"github.com/hkbarton/arango-cli/terminal"
)

type mainRunner struct{}

func (r mainRunner) Run(c *commands.Command, resultChan chan []string) {
	newState := map[string]string{
		"currentHost": "127.0.0.1",
		"currentDB":   "_system",
	}
	state.SetState(newState)
	resultChan <- nil
	close(resultChan)
}

func main() {
	entryCommand, err := commands.ParseCommand(os.Args)
	if err != nil {
		fmt.Println(err)
		return
	}
	runner := mainRunner{}
	entryResult := commands.RunCommand(entryCommand, runner)
	terminal.Output(entryResult)

	reader := bufio.NewReader(os.Stdin)
	for {
		commandString, _ := reader.ReadString('\n')
		if strings.TrimSpace(commandString) == "" {
			terminal.Output(nil)
			continue
		}
		command, err := commands.ParseCommand(strings.Split(commandString, " "))
		if err == nil {
			terminal.Output(commands.RunCommandByAction(command))
		} else {
			terminal.Output(nil)
		}
	}
}
