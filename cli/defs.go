package cli
import (
"flag"
"os"
"fmt"
)
var (
// CLI is the central struct characterizing the CLI
CLI = &Command{
description: []string{ "go-cli", "An example CLI app written using Go." },
flagSet: flag.CommandLine,
Subcommands: map[string]*Command{
"greet": {
description: []string{ "greet", "Get a greeting from your creation." },
flagSet: flag.NewFlagSet("greet", flag.ExitOnError),
},
},
}
name = CLI.Subcommands[`greet`].flagSet.String(`name`, `John Doe`, `Specify your name.`)
unexported = CLI.Subcommands[`greet`].flagSet.String(`unexported`, `Ohayou Gozaimasu!`, `This var won't be configurable through the config file.`)
customVar = CLI.Subcommands[`greet`].flagSet.Bool(`custom`, false, `Flag values can be stored in custom variables.`)
file = CLI.flagSet.String(`config`, `config.json`, `Path to configuration file.`)
gopher = CLI.flagSet.Bool(`gopher`, false, `Show something awesome.`)

appName = "go-cli"
appDesc = "An example CLI app written using Go."
appVersion = "1.2"

helpFlag = false
versionFlag = false
)
func init() {
args := os.Args
CLI.Subcommands[`greet`].flagSet.BoolVar(&helpFlag, `h`, false, `Display the help message.`)
CLI.Subcommands[`greet`].flagSet.BoolVar(&helpFlag, `help`, false, `Display the help message.`)
CLI.flagSet.BoolVar(&versionFlag, `version`, false, `Display version information.`)
CLI.flagSet.BoolVar(&helpFlag, `h`, false, `Display the help message.`)
CLI.flagSet.BoolVar(&helpFlag, `help`, false, `Display the help message.`)
CLI.parse(args[1:]...)
args = CLI.Args()
if len(args) > 0 {
switch args[0] {
case "greet":
CLI.Subcommands[`greet`].parse(args[1:]...)
default:
fmt.Printf("%s: invalid command\n", args[0])
CLI.usage(true)
os.Exit(2)
}
} else {
CLI.usage(false)
os.Exit(0)
}
}
