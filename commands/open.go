package commands

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hkbarton/arango-cli/state"

	"github.com/hkbarton/arango-cli/utils"
)

type OpenCommandRunner struct{}

func parseFilter(input string) map[string]string {
	if len(input) == 0 {
		return nil
	}
	result := make(map[string]string)
	filters := strings.Split(input, "&&")
	for _, filter := range filters {
		pair := strings.Split(filter, "=")
		key := strings.TrimSpace(pair[0])
		var value string
		if len(pair) > 1 {
			value = strings.TrimSpace(pair[1])
		}
		result[key] = value
	}
	return result
}

func parseFilterValue(input string) (value interface{}) {
	if input[0] == '"' && input[len(input)-1] == '"' {
		return strings.Trim(input, `"`)
	} else {
		boolreg, _ := regexp.Compile(`(?i)^(true|false)$`)
		intreg, _ := regexp.Compile(`^\d+$`)
		floatreg, _ := regexp.Compile(`^\d+\.\d+$`)
		switch {
		case boolreg.MatchString(input):
			value, _ = strconv.ParseBool(input)
		case intreg.MatchString(input):
			value, _ = strconv.Atoi(input)
		case floatreg.MatchString(input):
			value, _ = strconv.ParseFloat(input, 64)
		default:
			value = ""
		}
		return value
	}
}

func getCountOfCollection(collName string) int {
	col, err := state.CurrentDB().Collection(nil, collName)
	if err != nil {
		return 0
	}
	count, err := col.Count(nil)
	if err != nil {
		return 0
	}
	return int(count)
}

func (r OpenCommandRunner) Run(c *Command, resultChan chan interface{}) {
	defer close(resultChan)
	if c.Args == nil || len(c.Args) < 1 {
		resultChan <- errors.New("Open command needs collection target, e.g. open users")
		return
	}
	results := make([]string, 0, 30)
	target := c.Args[0]
	filterInput, hasFilter := c.Options["f"]
	var filterQuery string
	var filterArgs map[string]interface{}
	if hasFilter {
		filter := parseFilter(filterInput)
		filterQuery = "FILTER "
		filterArgs = make(map[string]interface{})
		for filterKey, filterValue := range filter {
			filterQuery = filterQuery + fmt.Sprintf("d.%s == @%s &&", filterKey, filterKey)
			filterArgs[filterKey] = parseFilterValue(filterValue)
		}
		filterQuery = strings.TrimSuffix(filterQuery, " &&")
	}
	queryResults, count, err := utils.QueryAll(nil, fmt.Sprintf(`
		FOR d IN %s
		%s
		LIMIT 10
		RETURN d
	`, target, filterQuery), filterArgs)
	if err != nil {
		resultChan <- err
		return
	}
	for _, item := range queryResults {
		if _, verbose := c.Options["v"]; verbose {
			results = append(results, item.JsonString())
		} else {
			itemStr := item.String()
			if len(itemStr) > 200 {
				itemStr = fmt.Sprintf("%s...", itemStr[:200])
			}
			results = append(results, itemStr)
		}
	}
	results = append(results, "\n")
	docCount := getCountOfCollection(target)
	results = append(results,
		fmt.Sprintf(utils.Color("Total doc count in this collection: ", utils.CyanColor)+"%d", docCount))
	if count > 0 && count < docCount {
		results = append(results,
			fmt.Sprintf("Showing first %s documents, you can do %s to view next page",
				utils.Color(strconv.FormatInt(int64(count), 10), utils.CyanColor),
				utils.Color("next", utils.CyanColor),
			),
		)
	}
	resultChan <- results
}
