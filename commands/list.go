package commands

import (
	"errors"
	"reflect"
	"sort"
	"strings"

	"github.com/hkbarton/arango-cli/state"
)

type ListCommandRunner struct{}

type objWithName interface {
	Name() string
}
type objGetter func() ([]objWithName, error)
type stringFilter func(input string) bool

func getObjectNames(objGetter objGetter, filter stringFilter) ([]string, error) {
	objs, err := objGetter()
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(objs))
	for _, obj := range objs {
		name := strings.TrimSpace(obj.Name())
		if len(name) > 0 && (filter == nil || filter(name)) {
			names = append(names, name)
		}
	}
	sort.Slice(names, func(i, j int) bool {
		return names[i] < names[j]
	})
	return names, nil
}

func convertObjGetterResult(objs interface{}, err error) ([]objWithName, error) {
	objSlice := reflect.ValueOf(objs)
	len := objSlice.Len()
	result := make([]objWithName, len)
	for i := 0; i < len; i++ {
		result[i] = objSlice.Index(i).Interface().(objWithName)
	}
	return result, err
}

func systemNameFilter(name string) bool {
	return !strings.HasPrefix(name, "_")
}

func (r ListCommandRunner) Run(c *Command, resultChan chan interface{}) {
	defer close(resultChan)
	if c.Args == nil || len(c.Args) < 1 {
		resultChan <- errors.New("List command needs arguments, e.g. list db")
		return
	}
	var err error
	var result interface{}
	target := c.Args[0]
	_, willShowAll := c.Options["a"]
	var filter stringFilter
	if !willShowAll {
		filter = systemNameFilter
	}
	switch target {
	case "db":
		result, err = getObjectNames(func() ([]objWithName, error) {
			return convertObjGetterResult(state.DBClient().Databases(nil))
		}, filter)
	case "coll":
		result, err = getObjectNames(func() ([]objWithName, error) {
			return convertObjGetterResult(state.CurrentDB().Collections(nil))
		}, filter)
	default:
		err = errors.New("Unknown list target " + target)
	}
	if err != nil {
		resultChan <- err
	} else {
		resultChan <- result
	}
}
