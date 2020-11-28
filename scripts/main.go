package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type (
	// Command is a type which characterizes the commands/subcommands
	command struct {
		Name        string    `yaml:"name"`
		Description string    `yaml:"description"`
		Flags       flags     `yaml:"flags"`
		Subcommands []command `yaml:"subcommands"`
	}

	flags []map[string]string
)

var (
	appVersion, cfgFlag string
	cli                 = command{}
	file, _             = ioutil.ReadFile("cli/schema.yaml")
)

func init() {
	yaml.Unmarshal(file, &cli)
	yaml.Unmarshal(file, &struct {
		Cfg *string `yaml:"configFlag"`
		Ver *string `yaml:"version"`
	}{&cfgFlag, &appVersion})
}

func main() {
	generateCLI()
	generateConfig()
}
