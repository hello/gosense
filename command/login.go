package command

import (
	"code.google.com/p/gopass"
	"fmt"
	"github.com/mitchellh/cli"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// LoginCommand is a Command implementation that attempts to
// login to sense api
type LoginCommand struct {
	Ui cli.Ui
}

func (c *LoginCommand) Help() string {
	helpText := `
Usage: sense login

  Retrieves a token from sense api and stores it locally
`
	return strings.TrimSpace(helpText)
}

func (c *LoginCommand) Run(args []string) int {
	if len(args) > 0 {
		c.Ui.Error("Please enter username and password at the prompt.")
		c.Ui.Error("")
		c.Ui.Error(c.Help())
		return 1
	}

	username, err := c.Ui.Ask("Username: ")

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading username from prompt: %s", err))
		return 1
	}

	password, err := gopass.GetPass("Password: ")

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading password from prompt: %s", err))
		return 1
	}

	client, err := AuthenticatedSenseClient(false)

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error connecting to sense API: %s", err))
		return 1
	}
	// defer client.Close()
	token, _, err := client.Tokens.Login(username, password)

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error retrieving token: %s", err))
		return 1
	}
	fmt.Printf("token = %s\n", token.Value)

	home := os.Getenv("HOME")
	settingsPath := filepath.Join(home, SettingsFileName)

	err = ioutil.WriteFile(settingsPath, []byte(token.Value), 0755)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error persisting token to file: %s", err))
		return 1
	}

	c.Ui.Info(fmt.Sprintf("Successfully logged in. Access token persisted to: ~/%s", SettingsFileName))
	return 0
}

func (c *LoginCommand) Synopsis() string {
	return "Log in to sense cloud"
}
