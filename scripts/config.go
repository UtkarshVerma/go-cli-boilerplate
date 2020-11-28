package main

import (
	"fmt"
	"os"

	"github.com/utkarshverma/go-cli-boilerplate/utils"
)

const configLayout = `package config
const configFlag = "%s"
type (
config struct {
File string` + "`json:\"-\"`" + `
%s}
)
`

var config string

func generateConfig() {
	cli.traverse()
	out, _ := os.Create("config/defs.go")
	out.Write([]byte(fmt.Sprintf(configLayout, cfgFlag, config)))
}

func (cmd command) traverse() {
	if cmd.Name != cli.Name {
		config += utils.ToPascalCase(cmd.Name) + " struct {\n"
	}

	if flags := cmd.Flags; len(flags) > 0 {
		for _, flag := range flags {
			name := flag["name"]
			if cmd.Name == cli.Name && name == cfgFlag {
				continue
			}
			name = utils.ToPascalCase(name)

			jsonTag := utils.ToSnakeCase(flag["name"]) + ",omitempty"
			if flag["unexported"] == "true" {
				jsonTag = "-"
			}

			tag := fmt.Sprintf("`json:\"%s\"`", jsonTag)
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
		config += fmt.Sprintf("} `json:\"%s,omitempty\"`\n", utils.ToSnakeCase(cmd.Name))
	}
}
