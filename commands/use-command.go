package commands

import (
	"errors"

	"github.com/hkbarton/arango-cli/state"
)

type UseCommandRunner struct{}

func (r UseCommandRunner) Run(c *Command, resultChan chan interface{}) {
	defer close(resultChan)
	if c.Args == nil || len(c.Args) < 1 {
		resultChan <- errors.New("Use command needs a database name as argument")
		return
	}
	db := c.Args[0]
	err := state.SetCurrentDB(db)
	if err != nil {
		resultChan <- err
	}
}
