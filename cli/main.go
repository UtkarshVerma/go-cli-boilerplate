package cli

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/utkarshverma/go-cli-boilerplate/utils"
)

type (
	// Command is a type which characterizes the commands/subcommands
	Command struct {
		name, description string
		flagSet           *flag.FlagSet

		Subcommands map[string]*Command
	}
)

// Args returns the arguments passed to `cmd` command/subcommand.
func (cmd *Command) Args() []string {
	return cmd.flagSet.Args()
}

// Arg returns the nth argument passed to `cmd` command/subcommand.
func (cmd *Command) Arg(n int) string {
	if n >= len(cmd.Args()) {
		return ""
	}
	return cmd.Args()[n]
}

// GetFlags returns the flags associated with `cmd` command/subcommand.
func (cmd *Command) GetFlags() *map[string]interface{} {
	flags := map[string]interface{}{}
	cmd.flagSet.VisitAll(
		func(f *flag.Flag) {
			flags[f.Name] = f.Value.(flag.Getter).Get()
		},
	)
	return &flags
}

// GetRaisedFlags returns the raised flags associated with `cmd` command/subcommand.
func (cmd *Command) GetRaisedFlags(parent *Command) (flags []string) {
	parentArgs := os.Args
	if parent != nil {
		parentArgs = parent.Args()
	}

	args := parentArgs[:len(parentArgs)-len(cmd.Args())]
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
		}
	}
	return
}

// GetFlag returns the value of `cmd` command/subcommand's `flag` flag.
func (cmd *Command) GetFlag(flag string) interface{} {
	flags := *cmd.GetFlags()
	return flags[flag]
}

// GetDefault returns the default value `cmd` command/subcommand's `f` flag.
func (cmd *Command) GetDefault(f string) (v interface{}) {
	val := cmd.flagSet.Lookup(f).DefValue
	valType := utils.TypeOf(val)

	v = val
	switch valType {
	case "bool":
		v, _ = strconv.ParseBool(val)
	case "int64":
		v, _ = strconv.ParseInt(val, 10, 64)
	case "float64":
		v, _ = strconv.ParseFloat(val, 64)
	}
	return
}

func (cmd *Command) parse(args ...string) {
	for _, subcmd := range cmd.Subcommands {
		if len(subcmd.Subcommands) > 0 {
			subcmd.describe()
		}
	}
	cmd.flagSet.Usage = func() { cmd.usage(false) }

	if cmd == App {
		flag.CommandLine.Parse(args)

		if versionFlag {
			fmt.Printf("%s v%s\n", App.name, appVersion)
			os.Exit(0)
		}
	} else {
		cmd.flagSet.Parse(args)
	}

	// WEIRD: Main plan was to use flag.ErrHelp but it isn't working as the docs say it should
	if helpFlag {
		cmd.flagSet.Usage()
		os.Exit(0)
	}
}

func (cmd *Command) usage(isInvalid bool) {
	name := cmd.flagSet.Name()
	if strings.Contains(name, App.name) {
		if !isInvalid {
			fmt.Println(App.description)
			fmt.Println()
		}
		fmt.Printf("Usage: %s COMMAND [OPTIONS...]\n", App.name)
	} else {
		args := App.Args()
		if isInvalid || helpFlag {
			args = args[:len(args)-1]
		}

		var usage string = App.name
		if len(args) > 0 {
			usage += " " + strings.Join(args, " ")
		}
		if len(cmd.Subcommands) > 0 {
			usage += " COMMAND"
		}

		fmt.Printf("Usage: %s [OPTIONS...]\n", usage)
	}

	fmt.Printf("For help: %s [COMMAND] -help\n", App.name)

	if len(cmd.Subcommands) > 0 {
		fmt.Println("\nAvailable commands:")

		keys := make([]string, len(cmd.Subcommands))
		i := 0
		for key := range cmd.Subcommands {
			keys[i] = key
			i++
		}
		sort.Strings(keys)
		for _, key := range keys {
			subcmd := cmd.Subcommands[key]
			fmt.Printf("\t%-20s%s\n", subcmd.name, subcmd.description)
		}
	}

	flags := cmd.flagSet
	nFlags := 0
	flags.VisitAll(func(f *flag.Flag) { nFlags++ })
	if nFlags > 2 {
		fmt.Println("\nAvailable options:")
		flags.VisitAll(
			func(f *flag.Flag) {
				if f.Name != "help" && f.Name != "h" {
					if _, ok := f.Value.(flag.Getter).Get().(string); ok {
						f.DefValue = fmt.Sprintf("\"%s\"", f.DefValue)
					}
					fmt.Printf("\t-%-19s%s\n", f.Name, f.Usage)
					if f.Name != "version" && f.DefValue != `""` && f.DefValue != "false" {
						fmt.Printf("\t %-21s(default: %v)\n", " ", f.DefValue)
					}
				}
			},
		)
	}
}

func (cmd *Command) describe() {
	subcmds := make([]string, len(cmd.Subcommands))
	i := 0
	for _, v := range cmd.Subcommands {
		subcmds[i] = "'" + v.flagSet.Name() + "'"
		i++
	}
	sort.Strings(subcmds)

	l := len(subcmds)
	desc := []string{fmt.Sprintf("%s SUBCOMMAND", cmd.flagSet.Name())}
	if l > 1 {
		desc = append(desc, fmt.Sprintf("SUBCOMMAND can be %s or %s.",
			strings.Join(subcmds[:l-1], ", "),
			subcmds[l-1],
		))
	} else {
		desc = append(desc, fmt.Sprintf("SUBCOMMAND can be %s.", subcmds[0]))
	}
	cmd.name, cmd.description = desc[0], desc[1]
}
