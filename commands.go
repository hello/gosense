package main

import (
	"github.com/hello/sense/command"
	"github.com/mitchellh/cli"
	"os"
)

// Commands is the mapping of all the available Spark commands.
var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.ColoredUi{
		InfoColor:  cli.UiColorGreen,
		ErrorColor: cli.UiColorRed,
		Ui: &cli.BasicUi{
			Writer: os.Stdout,
			Reader: os.Stdin,
		},
	}

	Commands = map[string]cli.CommandFactory{
		"login": func() (cli.Command, error) {
			return &command.LoginCommand{
				Ui: ui,
			}, nil
		},
		"whoami": func() (cli.Command, error) {
			return &command.WhoAmICommand{
				Ui: ui,
			}, nil
		},
		"register": func() (cli.Command, error) {
			return &command.RegisterCommand{
				Ui: ui,
			}, nil
		},
		"update": func() (cli.Command, error) {
			return &command.UpdateAccountCommand{
				Ui: ui,
			}, nil
		},

		"upload": func() (cli.Command, error) {
			return &command.UploadDataCommand{
				Ui: ui,
			}, nil
		},
		"version": func() (cli.Command, error) {
			return &command.VersionCommand{
				Version: Version,
				Ui:      ui,
			}, nil
		},
	}
}
