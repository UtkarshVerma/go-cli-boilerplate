package cli
import (
"flag"
"os"
"fmt"
)
var (
// App is the central struct characterizing the CLI
App = &Command{
name: "go-cli",
description: "An example CLI app written using Go.",
flagSet: flag.CommandLine,
Subcommands: map[string]*Command{
"greet": {
name: "greet",
description: "Get a greeting from your creation.",
flagSet: flag.NewFlagSet("greet", flag.ExitOnError),
},
},
}
name = App.Subcommands[`greet`].flagSet.String(`name`, `John Doe`, `Specify your name.`)
unexported = App.Subcommands[`greet`].flagSet.String(`unexported`, `Ohayou Gozaimasu!`, `This var won't be configurable through the config file.`)
customVar = App.Subcommands[`greet`].flagSet.Bool(`custom`, false, `Flag values can be stored in custom variables.`)
config = App.flagSet.String(`config`, `config.json`, `Path to configuration file.`)
gopher = App.flagSet.Bool(`gopher`, false, `Show something awesome.`)
appVersion = "0.1"

helpFlag = false
versionFlag = false
)
func init() {
args := os.Args
App.flagSet.BoolVar(&versionFlag, `version`, false, `Display version information.`)
App.Subcommands[`greet`].flagSet.BoolVar(&helpFlag, `h`, false, `Display the help message.`)
App.Subcommands[`greet`].flagSet.BoolVar(&helpFlag, `help`, false, `Display the help message.`)
App.flagSet.BoolVar(&helpFlag, `h`, false, `Display the help message.`)
App.flagSet.BoolVar(&helpFlag, `help`, false, `Display the help message.`)
App.parse(args[1:]...)
args = App.Args()
if len(args) > 0 {
switch args[0] {
case "greet":
App.Subcommands[`greet`].parse(args[1:]...)
default:
fmt.Printf("%s: invalid command\n\n", args[0])
App.usage(true)
os.Exit(2)
}
} else {
App.usage(false)
os.Exit(0)
}
}
