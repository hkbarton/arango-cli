package commands

import (
	"errors"
	"strings"

	"github.com/hkbarton/arango-cli/state"
)

type ListCommandRunner struct{}

func (r ListCommandRunner) Run(c *Command, resultChan chan interface{}) {
	if c.Args == nil || len(c.Args) < 1 {
		resultChan <- errors.New("List command need arguments, e.g. list db")
		return
	}
	target := c.Args[0]
	switch target {
	case "db":
		databases, err := state.DBClient().Databases(nil)
		if err != nil {
			resultChan <- err
			break
		}
		dbNames := make([]string, 0, len(databases))
		for _, db := range databases {
			name := db.Name()
			if len(strings.TrimSpace(name)) > 0 {
				dbNames = append(dbNames, name)
			}
		}
		resultChan <- dbNames
	default:
		resultChan <- errors.New("Unknown list target " + target)
	}
	close(resultChan)
}
