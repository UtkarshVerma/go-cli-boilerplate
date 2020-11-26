package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/utkarshverma/go-cli-boilerplate/utils"
)

const cliLayout = `package cli
import (
"flag"
"os"
"fmt"
)
var (
%s
helpFlag = false
)
func init() {
args := os.Args
%s}
`

var cliVar, cliStruct, initFunc, flagStr string

func generateCLI() {
	// TODO: Think of a better value for `cliVar`, the one below isn't modular
	// cliVar = utils.ToPascalCase(cli.Name)
	cliVar = "CLI"
	cli.define()
	cli.defineInitFunc()

	out, _ := os.Create("cli/defs.go")
	out.Write([]byte(fmt.Sprintf(cliLayout, cliStruct+flagStr, initFunc)))
}

func (cmd command) define(cmds ...string) {
	nestLevel := len(cmds)
	if nestLevel == 0 {
		cliStruct = fmt.Sprintf(
			"// %s is the central struct characterizing the CLI\n"+
				"%s = &Command{\n"+
				"description: []string{ \"%s\", \"%s\" },\n"+
				"flagSet: flag.CommandLine,\n",
			cliVar, cliVar,
			cmd.Name, cmd.Description,
		)
	} else {
		cliStruct += fmt.Sprintf(
			"\"%s\": {\n"+
				"description: []string{ \"%s\", \"%s\" },\n"+
				"flagSet: flag.NewFlagSet(\"%s\", flag.ExitOnError),\n",
			cmd.Name,
			cmd.Name, cmd.Description,
			cmd.Name,
		)
	}

	// Define subcommands
	for i := range cmd.Subcommands {
		if i == 0 {
			cliStruct += "Subcommands: map[string]*Command{\n"
		}
		cmd.Subcommands[i].define(append(cmds, cmd.Subcommands[i].Name)...)
		if i == len(cmd.Subcommands)-1 {
			cliStruct += "},\n"
		}
	}

	if nestLevel == 0 {
		cliStruct += "}\n"
	} else {
		cliStruct += "},\n"
	}

	cmd.defineFlags(cmds...)
}

func (cmd command) defineFlags(cmds ...string) {
	flagSet := cliVar
	for _, c := range cmds {
		flagSet += ".Subcommands[`" + c + "`]"
	}
	for _, flag := range cmd.Flags {
		value := flag["default"]
		kind := strings.Title(utils.TypeOf(value))
		if kind == "String" {
			value = "`" + value + "`"
		}

		name := flag["name"]
		if v, ok := flag["var"]; ok {
			name = v
		}
		name = utils.ToCamelCase(name)

		flagStr += fmt.Sprintf(
			"%s = %s.flagSet.%s(`%s`, %s, `%s`)\n",
			name, flagSet, kind, flag["name"], value, flag["description"],
		)
	}
	initFunc += flagSet + ".flagSet.BoolVar(&helpFlag, `h`, false, `Display the help message.`)\n"
	initFunc += flagSet + ".flagSet.BoolVar(&helpFlag, `help`, false, `Display the help message.`)\n"
}

func (cmd command) defineInitFunc(cmds ...string) {
	flagSet := cliVar
	for _, c := range cmds {
		flagSet += ".Subcommands[`" + c + "`]"
	}

	initFunc += flagSet + ".parse(args[1:]...)\n"
	if len(cmd.Subcommands) > 0 {
		initFunc += "args = " + flagSet + ".Args()\n" +
			"if len(args) > 0 {\n" +
			"switch args[0] {\n"

		for i, v := range cmd.Subcommands {
			initFunc += "case \"" + v.Name + "\":\n"
			cmd.Subcommands[i].defineInitFunc(append(cmds, v.Name)...)
		}

		initFunc += "default:\n" +
			`fmt.Printf("%s: invalid command\n", args[0])` + "\n" +
			flagSet + ".usage(true)\n" +
			"os.Exit(2)\n" +
			"}\n} else {\n" +
			flagSet + ".usage(false)\n" +
			"os.Exit(0)\n" +
			"}\n"
	}
}
