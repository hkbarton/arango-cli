package commands

import (
	"errors"
	"fmt"

	"github.com/hkbarton/arango-cli/state"
)

type OpenCommandRunner struct{}

func (r OpenCommandRunner) Run(c *Command, resultChan chan interface{}) {
	defer close(resultChan)
	if c.Args == nil || len(c.Args) < 1 {
		resultChan <- errors.New("Open command needs collection target, e.g. open users")
		return
	}
	target := c.Args[0]
	col, err := state.CurrentDB().Collection(nil, target)
	if err != nil {
		resultChan <- err
		return
	}
	results := make([]string, 0, 11)
	count, err := col.Count(nil)
	if err != nil {
		resultChan <- err
		return
	}
	results = append(results, fmt.Sprintf("Total count: %d", count))
	resultChan <- results
}
