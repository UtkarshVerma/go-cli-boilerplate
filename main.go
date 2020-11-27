package main

import (
	"fmt"
	"strings"

	"github.com/utkarshverma/go-cli-boilerplate/cli"
	"github.com/utkarshverma/go-cli-boilerplate/config"
)

var (
	conf = config.Config
)

//go:generate go run ./scripts

func main() {
	switch cli.App.Arg(0) {
	case "greet":
		name := cli.App.Subcommands["greet"].GetFlag("name").(string)
		firstName := strings.Split(name, " ")[0]
		fmt.Printf("Ohayou gozaimasu, %s-san.\n", firstName)
	}
}
