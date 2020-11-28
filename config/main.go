package config

import (
	"log"
	"reflect"

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

func (config *config) update() {
	update(cli.App, reflect.ValueOf(config).Elem())
}

func update(from *cli.Command, to reflect.Value) {
	fields := to.Type()
	for i := 0; i < fields.NumField(); i++ {
		// Reconstruct the flag name, and don't update config flag
		field := utils.ToKebabCase(fields.Field(i).Name)
		if field == "file" {
			continue
		}

		copyTo := to.Field(i)
		if to.Field(i).Kind() == reflect.Struct {
			update(from.Subcommands[field], copyTo)
		} else {
			mustSet := from.GetDefault(field) != from.GetFlag(field)

			switch copyTo.Kind() {
			case reflect.String:
				if copyTo.CanSet() && mustSet {
					copyTo.SetString(from.GetFlag(field).(string))
				}
			case reflect.Int:
				if copyTo.CanSet() && mustSet {
					copyTo.SetInt(int64(from.GetFlag(field).(int)))
				}
			case reflect.Bool:
				if copyTo.CanSet() && mustSet {
					copyTo.SetBool(from.GetFlag(field).(bool))
				}
			}
		}
	}
}
