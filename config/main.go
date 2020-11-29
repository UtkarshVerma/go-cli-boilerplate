package config

import (
	"log"
	"reflect"
	"strings"

	"github.com/utkarshverma/go-cli-boilerplate/cli"
	"github.com/utkarshverma/go-cli-boilerplate/utils"
)

var (
	// App is the central configuration struct.
	App = &config{}
)

func init() {
	if configFlag != "" {
		App.ConfigFile = cli.App.GetFlag(configFlag).(string)
		App.init()

		// Create config file if not present
		if utils.FileExists(App.ConfigFile) {
			err := utils.ReadJSON(App.ConfigFile, App)
			if err != nil {
				log.Fatal(err)
			}
			App.update()
		} else {
			App.update()
			err := utils.WriteJSON(App, App.ConfigFile)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (config *config) init() {
	traverse(nil, cli.App, reflect.ValueOf(config).Elem(), true)
}

func (config *config) update() {
	traverse(nil, cli.App, reflect.ValueOf(config).Elem(), false)
}

func traverse(parent *cli.Command, from *cli.Command, to reflect.Value, mustInit bool) {
	fields := to.Type()
	for i := 0; i < fields.NumField(); i++ {
		// Reconstruct the flag name, and don't update config flag
		field := fields.Field(i).Name
		if field == "ConfigFile" {
			continue
		}
		field = utils.ToKebabCase(field)

		copyTo := to.Field(i)
		if to.Field(i).Kind() == reflect.Struct {
			traverse(from, from.Subcommands[field], copyTo, mustInit)
		} else {
			getVal := from.GetDefault
			mustSet := true
			if !mustInit {
				getVal = from.GetFlag
				mustSet = false

				raisedFlags := from.GetRaisedFlags(parent)
				for _, flag := range raisedFlags {
					if strings.Contains(flag, field) {
						mustSet = true
						break
					}
				}
			}

			if copyTo.CanSet() && mustSet {
				switch copyTo.Kind() {
				case reflect.String:
					copyTo.SetString(getVal(field).(string))
				case reflect.Bool:
					copyTo.SetBool(getVal(field).(bool))
				case reflect.Float64:
					copyTo.SetFloat(getVal(field).(float64))
				case reflect.Int64:
					copyTo.SetInt(getVal(field).(int64))
				}
			}
		}
	}
}
