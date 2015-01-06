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
		c.Ui.Error(fmt.Sprintf("Error converting to int: %s", err))
		return 1
	}

	deviceId, err := c.Ui.Ask("Device id:\n ex: 97E5E9F490AA306B (jackson)")
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading device_id: %s", err))
		return 1
	}

	aesKey, err := c.Ui.Ask("AES key(hex):\n ex: 199C126323ABC61CCB291BECC337FB66 (jackson)")
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading aes_key: %s", err))
		return 1
	}

	resp, err := client.Upload.Upload(int32(t), deviceId, aesKey)

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error sending/receiving data: %s", err))
		return 1
	}

	// c.Ui.Info(fmt.Sprintf("[OK] : %v", resp))
	c.Ui.Info(fmt.Sprintf("[-->] : batch size %d", resp.GetBatchSize()))
	c.Ui.Info(fmt.Sprintf("[-->] : audio_control %s", resp.GetAudioControl().GetAudioCaptureAction().Enum()))
	c.Ui.Info("[-->] : alarm:")
	c.Ui.Info(fmt.Sprintf("\tringtone_id %d", resp.GetAlarm().GetRingtoneId()))
	c.Ui.Info(fmt.Sprintf("\tring_offset_from_now_in_second %d", resp.GetAlarm().GetRingOffsetFromNowInSecond()))
	c.Ui.Info(fmt.Sprintf("\tstart_time %d", resp.GetAlarm().GetStartTime()))
	c.Ui.Info(fmt.Sprintf("\tend_time %d", resp.GetAlarm().GetEndTime()))
	c.Ui.Info(fmt.Sprintf("[-->] %d files to download:", len(resp.GetFiles())))
	for _, file := range resp.GetFiles() {
		c.Ui.Info(fmt.Sprintf("\tfile: %s", file.GetUrl()))
	}
	return 0
}

func (c *UploadDataCommand) Synopsis() string {
	return "upload periodic data to cloud"
}
