package command

import (
	"fmt"
	"github.com/mitchellh/cli"
	"strconv"
	"strings"
)

// LogoutCommand is a Command implementation that attempts to
// delete the access token stored locally
type UploadDataCommand struct {
	Ui cli.Ui
}

func (c *UploadDataCommand) Help() string {
	helpText := `
Usage: sense whoami

  Empties local settings file which contains the most recent access token.
`
	return strings.TrimSpace(helpText)
}

func (c *UploadDataCommand) Run(args []string) int {
	if len(args) > 0 {
		c.Ui.Error("This command does not accept any argument.")
		c.Ui.Error("")
		c.Ui.Error(c.Help())
		return 1
	}

	client, err := SenseProtobufClient()

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed connecting to sense cloud: %s", err))
		return 1
	}

	temp, err := c.Ui.Ask("Temp: ")

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading name from prompt: %s", err))
		return 1
	}

	t, err := strconv.Atoi(temp)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading name from prompt: %s", err))
		return 1
	}
	for i := 0; i < 100; i++ {
		resp, err := client.Upload.Upload(int32(t))

		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error uploading data: %s", err))
			return 1
		}
		c.Ui.Info(fmt.Sprintf("[OK] : %v", resp))
	}

	return 0
}

func (c *UploadDataCommand) Synopsis() string {
	return "upload periodic data to cloud"
}
