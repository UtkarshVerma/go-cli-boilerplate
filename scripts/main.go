package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type (
	// Command is a type which characterizes the commands/subcommands
	command struct {
		Name        string    `yaml:"name,omitempty"`
		Description string    `yaml:"description,omitempty"`
		Var         string    `yaml:"var,omitempty"`
		Flags       flags     `yaml:"flags,omitempty"`
		Subcommands []command `yaml:"subcommands,omitempty"`
	}

	flags []map[string]string
)

var (
	cli     = command{}
	file, _ = ioutil.ReadFile("cli/schema.yaml")
)

func init() {
	yaml.Unmarshal(file, &cli)
}

func main() {
	generateCLI()
	generateConfig()
}
