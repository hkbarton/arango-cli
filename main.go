package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hkbarton/arango-cli/commands"
	"github.com/hkbarton/arango-cli/state"
	"github.com/hkbarton/arango-cli/utils"

	driver "github.com/arangodb/go-driver"
	driverHttp "github.com/arangodb/go-driver/http"
)

type mainRunner struct{}

// Default DB information
const (
	DefaultDBHost = "127.0.0.1"
	DefaultDBPort = 8529
)

func (r mainRunner) Run(c *commands.Command, resultChan chan interface{}) {
	defer close(resultChan)
	host, hostExists := c.Options["h"]
	port, portErr := strconv.ParseInt(c.Options["p"], 10, 8)
	if !hostExists {
		host = DefaultDBHost
	}
	if portErr != nil {
		port = DefaultDBPort
	}

	connErr := func() {
		err := utils.FatalError{Message: fmt.Sprintf("Failed to connect %s:%d", host, port)}
		resultChan <- err
	}

	conn, err := driverHttp.NewConnection(driverHttp.ConnectionConfig{
		Endpoints: []string{fmt.Sprintf("http://%s:%d", host, port)},
	})
	if err != nil {
		connErr()
		return
	}
	client, err := driver.NewClient(driver.ClientConfig{
		Connection: conn,
	})
	if err != nil {
		connErr()
		return
	}

	state.SetState(map[string]interface{}{
		"currentHost": host + ":" + strconv.Itoa(int(port)),
		"dbClient":    client,
	})

	err = state.SetCurrentDB("_system")
	if err != nil {
		connErr()
		return
	}
}

func main() {
	entryCommand, err := commands.Parse(os.Args)
	if err != nil {
		fmt.Println(err)
		return
	}
	runner := mainRunner{}
	entryResult := commands.Run(entryCommand, runner)
	utils.Output(entryResult)

	reader := bufio.NewReader(os.Stdin)
	for {
		commandString, _ := reader.ReadString('\n')
		if strings.TrimSpace(commandString) == "" {
			utils.Output(nil)
			continue
		}
		command, err := commands.Parse(strings.Split(commandString, " "))
		if err == nil {
			utils.Output(commands.RunByAction(command))
		} else {
			utils.Output(nil)
		}
	}
}
