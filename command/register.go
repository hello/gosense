package command

import (
	"fmt"
	"github.com/hello/sense/sense"
	"github.com/mitchellh/cli"
	"strings"
)

// LogoutCommand is a Command implementation that attempts to
// delete the access token stored locally
type RegisterCommand struct {
	Ui cli.Ui
}

func (c *RegisterCommand) Help() string {
	helpText := `
Usage: sense register

  Empties local settings file which contains the most recent access token.
`
	return strings.TrimSpace(helpText)
}

func (c *RegisterCommand) Run(args []string) int {
	if len(args) > 0 {
		c.Ui.Error("This command does not accept any argument.")
		c.Ui.Error("")
		c.Ui.Error(c.Help())
		return 1
	}

	email, err := c.Ui.Ask("Email: ")

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading email from prompt: %s", err))
		return 1
	}

	name, err := c.Ui.Ask("Name: ")

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading name from prompt: %s", err))
		return 1
	}

	password, err := c.Ui.Ask("Password: ")

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading password from prompt: %s", err))
		return 1
	}

	client, err := AuthenticatedSenseClient(false)

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed connecting to sense API: %s", err))
		return 1
	}

	reg := sense.NewRegistration(name, email, password)

	account, _, err := client.Account.Register(reg)

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error registering account info: %s", err))
		return 1
	}

	c.Ui.Info(fmt.Sprintf("Hello %s (%v)", account.Name, account))
	return 0
}

func (c *RegisterCommand) Synopsis() string {
	return "register new account"
}
