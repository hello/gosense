package command

import (
	"fmt"
	"github.com/mitchellh/cli"
	"strings"
)

// LogoutCommand is a Command implementation that attempts to
// delete the access token stored locally
type UpdateAccountCommand struct {
	Ui cli.Ui
}

func (c *UpdateAccountCommand) Help() string {
	helpText := `
Usage: sense whoami

  Empties local settings file which contains the most recent access token.
`
	return strings.TrimSpace(helpText)
}

func (c *UpdateAccountCommand) Run(args []string) int {
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

	account, _, err := client.Account.Me()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed getting account from server: %s", err))
		return 1
	}

	c.Ui.Info(fmt.Sprintf("Account last updated: %d", account.LastModified))

	name, err := c.Ui.Ask(fmt.Sprintf("Name (was %s): ", account.Name))

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading name from prompt: %s", err))
		return 1
	}

	dob, err := c.Ui.Ask(fmt.Sprintf("DOB (was %s): ", account.DOB))

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading DOB from prompt: %s", err))
		return 1
	}

	account.Name = name
	account.DOB = dob

	updated_account, _, err := client.Account.Update(&account)

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error getting account info: %s", err))
		return 1
	}

	c.Ui.Info(fmt.Sprintf("Hello %s (%v)", updated_account, account))
	return 0
}

func (c *UpdateAccountCommand) Synopsis() string {
	return "Update account command"
}
