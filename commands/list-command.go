package commands

import (
	"errors"
	"reflect"
	"strings"

	"github.com/hkbarton/arango-cli/state"
)

type ListCommandRunner struct{}

type objWithName interface {
	Name() string
}
type objGetter func() ([]objWithName, error)

func getObjectNames(objGetter objGetter) ([]string, error) {
	objs, err := objGetter()
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(objs))
	for _, obj := range objs {
		name := obj.Name()
		if len(strings.TrimSpace(name)) > 0 {
			names = append(names, name)
		}
	}
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

func (r ListCommandRunner) Run(c *Command, resultChan chan interface{}) {
	defer close(resultChan)
	if c.Args == nil || len(c.Args) < 1 {
		resultChan <- errors.New("List command needs arguments, e.g. list db")
		return
	}
	var err error
	var result interface{}
	target := c.Args[0]
	switch target {
	case "db":
		result, err = getObjectNames(func() ([]objWithName, error) {
			return convertObjGetterResult(state.DBClient().Databases(nil))
		})
	case "coll":
		result, err = getObjectNames(func() ([]objWithName, error) {
			return convertObjGetterResult(state.CurrentDB().Collections(nil))
		})
	default:
		err = errors.New("Unknown list target " + target)
	}
	if err != nil {
		resultChan <- err
	} else {
		resultChan <- result
	}
}
