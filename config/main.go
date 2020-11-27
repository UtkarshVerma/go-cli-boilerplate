package config

import (
	"log"
	"reflect"

	"github.com/utkarshverma/go-cli-boilerplate/cli"
	"github.com/utkarshverma/go-cli-boilerplate/utils"
)

var (
	// TimeLayout is the time layout which is followed everywhere
	TimeLayout = "2006-01-02T15:04:05-07:00"

	// Config is the central configuration struct
	Config = &config{
		File: cli.CLI.GetFlag("config").(string),
	}
)

func init() {
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

func (config *config) update() {
	update(cli.CLI, reflect.ValueOf(config).Elem())
}

func update(from *cli.Command, to reflect.Value) {
	fields := to.Type()
	for i := 0; i < fields.NumField(); i++ {
		// Reconstruct the flag name
		name := fields.Field(i).Name
		if nameTag := fields.Field(i).Tag.Get("name"); nameTag != "" {
			name = nameTag
		}
		name = utils.ToKebabCase(name)

		copyTo := to.Field(i)
		if to.Field(i).Kind() == reflect.Struct {
			update(from.Subcommands[name], copyTo)
		} else {
			mustSet := copyTo.IsZero()

			switch copyTo.Kind() {
			case reflect.String:
				if copyTo.CanSet() && mustSet {
					copyTo.SetString(from.GetFlag(name).(string))
				}
			case reflect.Int:
				if copyTo.CanSet() && mustSet {
					copyTo.SetInt(int64(from.GetFlag(name).(int)))
				}
			case reflect.Bool:
				if copyTo.CanSet() && mustSet {
					copyTo.SetBool(from.GetFlag(name).(bool))
				}
			}
		}
	}
}
