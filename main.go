package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hkbarton/arango-cli/commands"
	"github.com/hkbarton/arango-cli/state"
	"github.com/hkbarton/arango-cli/terminal"

	driver "github.com/arangodb/go-driver"
	driverHttp "github.com/arangodb/go-driver/http"
)

type mainRunner struct{}

// Default DB information
const (
	DefaultDBHost = "127.0.0.1"
	DefaultDBPort = 8529
)

func (r mainRunner) Run(c *commands.Command, resultChan chan []string) {
	host, hostExists := c.Options["h"]
	port, portErr := strconv.ParseInt(c.Options["p"], 10, 8)
	if !hostExists {
		host = DefaultDBHost
	}
	if portErr != nil {
		port = DefaultDBPort
	}

	connErr := func() {
		resultChan <- []string{fmt.Sprintf("Failed to connect %s:%d", host, port)}
		close(resultChan)
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
	_, err = client.Database(nil, "_system")
	if err != nil {
		connErr()
		return
	}

	state.SetState(map[string]interface{}{
		"currentHost": host + ":" + strconv.Itoa(int(port)),
		"currentDB":   "_system",
		"dbClient":    client,
	})
	close(resultChan)
}

func main() {
	entryCommand, err := commands.Parse(os.Args)
	if err != nil {
		fmt.Println(err)
		return
	}
	runner := mainRunner{}
	entryResult := commands.Run(entryCommand, runner)
	terminal.Output(entryResult)

	reader := bufio.NewReader(os.Stdin)
	for {
		commandString, _ := reader.ReadString('\n')
		if strings.TrimSpace(commandString) == "" {
			terminal.Output(nil)
			continue
		}
		command, err := commands.Parse(strings.Split(commandString, " "))
		if err == nil {
			terminal.Output(commands.RunByAction(command))
		} else {
			terminal.Output(nil)
		}
	}
}
