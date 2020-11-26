package main

import (
	"fmt"
	"os"

	"github.com/utkarshverma/go-cli-boilerplate/utils"
)

const configLayout = `package config
type (
config struct {
%s}
)
`

var config string

func generateConfig() {
	cli.traverse()
	out, _ := os.Create("config/defs.go")
	out.Write([]byte(fmt.Sprintf(configLayout, config)))
}

func (cmd command) traverse() {
	if cmd.Name != cli.Name {
		name := cmd.Name
		if cmd.Var != "" {
			name = cmd.Var
		}
		name = utils.ToPascalCase(name)
		config += name + " struct {\n"
	}

	if flags := cmd.Flags; len(flags) > 0 {
		for _, flag := range flags {
			var isCustomVar bool

			name := flag["name"]
			if v, ok := flag["var"]; ok {
				name = v
				isCustomVar = true
			}
			name = utils.ToPascalCase(name)

			jsonTag := utils.ToSnakeCase(flag["name"]) + ",omitempty"
			if flag["unexported"] == "true" {
				jsonTag = "-"
			}

			tag := fmt.Sprintf("`json:\"%s\"", jsonTag)
			if isCustomVar {
				tag += fmt.Sprintf(" name:\"%s\",", flag["name"])
			}
			tag += "`"

			kind := utils.TypeOf(flag["default"])
			config += fmt.Sprintf("%s %s %s\n", name, kind, tag)
		}
	}

	if len(cmd.Subcommands) > 0 {
		for _, subcmd := range cmd.Subcommands {
			subcmd.traverse()
		}
	}

	if cmd.Name != cli.Name {
		tag := fmt.Sprintf("`json:\"%s,omitempty\"", utils.ToSnakeCase(cmd.Name))
		if cmd.Var != "" {
			tag += fmt.Sprintf(" name:\"%s\"", cmd.Name)
		}
		tag += "`"

		config += fmt.Sprintf("} %s\n", tag)
	}
}
