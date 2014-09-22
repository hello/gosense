package command

import (
	"fmt"
	"github.com/mitchellh/cli"
)

// VersionCommand is a Command implementation which prints the version.
type VersionCommand struct {
	Version string
	Ui      cli.Ui
}

func (c *VersionCommand) Help() string {
	return ""
}

func (c *VersionCommand) Run(_ []string) int {
	c.Ui.Output(fmt.Sprintf("sense v%s", c.Version))
	c.Ui.Output("Using sense API v1")
	return 0
}

func (c *VersionCommand) Synopsis() string {
	return "Prints the sense cli version"
}
