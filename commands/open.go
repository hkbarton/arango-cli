package commands

import (
	"errors"
	"fmt"

	"github.com/hkbarton/arango-cli/utils"
)

type OpenCommandRunner struct{}

func (r OpenCommandRunner) Run(c *Command, resultChan chan interface{}) {
	defer close(resultChan)
	if c.Args == nil || len(c.Args) < 1 {
		resultChan <- errors.New("Open command needs collection target, e.g. open users")
		return
	}
	target := c.Args[0]
	results := make([]string, 0, 12)
	queryResults, count, err := utils.QueryAll(nil, fmt.Sprintf(`
		FOR d IN %s
		LIMIT 10
		RETURN d
	`, target), nil)
	if err != nil {
		resultChan <- err
		return
	}
	for _, item := range queryResults {
		itemStr := item.String()
		if len(itemStr) > 100 {
			itemStr = fmt.Sprintf("%s...", itemStr[:100])
		}
		results = append(results, itemStr)
	}
	results = append(results, "\n")
	results = append(results, fmt.Sprintf(utils.Color("Total count: ", utils.CyanColor)+"%d", count))
	resultChan <- results
}
