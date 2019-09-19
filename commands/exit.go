package commands

import "os"

// ExitCommandRunner run exit command
type ExitCommandRunner struct{}

// Run run exit command
func (r ExitCommandRunner) Run(c *Command, resultChan chan interface{}) {
	os.Exit(0)
}
