package main

import (
	"fmt"
	"os/user"
	"strings"
	"time"

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
		greet()
	}
}

func greet() {
	name := config.Config.Greet.Name
	if name == "" {
		user, _ := user.Current()
		name = user.Name
	}
	firstName := strings.Split(name, " ")[0]

	var greeting string
	switch t := time.Now(); {
	case t.Hour() < 12:
		greeting = "Ohayou-gozaimasu"
	case t.Hour() < 16:
		greeting = "Konnichiwa"
	case t.Hour() < 20:
		greeting = "Konbanwa"
	default:
		greeting = "Oyasuminasai"
	}

	fmt.Printf("%s, %s-san!\n", greeting, firstName)
}
