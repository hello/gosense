package command

import (
	"fmt"
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

	client, err := AuthenticatedSenseClient(true)

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed connecting to sense API: %s", err))
		return 1
	}

	account, _, err := client.Account.Register()

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
