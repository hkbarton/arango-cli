package commands

import (
	"fmt"

	"github.com/hkbarton/arango-cli/state"
	"github.com/hkbarton/arango-cli/utils"
)

type InfoCommandRunner struct{}

func (r InfoCommandRunner) Run(c *Command, resultChan chan interface{}) {
	defer close(resultChan)
	systemDB, err := state.DBClient().Database(nil, "_system")
	if err != nil {
		resultChan <- err
		return
	}
	infoQueryResult, err := utils.QueryOne(systemDB, `
		FOR s IN _statistics
		SORT s._key DESC
		LIMIT 1
		RETURN {
			_key: s._key,
			cnt: s.client.httpConnections
		}
	`, nil)
	if err != nil {
		resultChan <- err
		return
	}
	results := make([]string, 0, 20)
	results = append(results, fmt.Sprintf(utils.Color("Total connections: ", utils.CyanColor)+"%v", infoQueryResult.Document["cnt"]))
	resultChan <- results
}
