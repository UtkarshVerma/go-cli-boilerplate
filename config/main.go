package config

import (
	"log"
	"reflect"
	"strings"

	"github.com/utkarshverma/go-cli-boilerplate/cli"
	"github.com/utkarshverma/go-cli-boilerplate/utils"
)

var (
	// Config is the central configuration struct.
	Config = &config{}
)

func init() {
	if configFlag != "" {
		Config.File = cli.App.GetFlag(configFlag).(string)
		Config.init()

		// Create config file if not present
		if utils.FileExists(Config.File) {
			err := utils.ReadJSON(Config.File, Config)
			if err != nil {
				log.Fatal(err)
			}
			Config.update()
		} else {
			Config.update()
			err := utils.WriteJSON(Config, Config.File)
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
		field := utils.ToKebabCase(fields.Field(i).Name)
		if field == "file" {
			continue
		}

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
