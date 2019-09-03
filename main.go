package main

import (
	"fmt"
	"os"

	"github.com/hkbarton/arango-cli/commands"
)

func main() {
	entryCommand, err := commands.ParseCommand(os.Args)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(entryCommand)
}
