package command

import (
	"errors"
	"github.com/hello/sense/sense"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const (
	SettingsFileName = ".sense"
)

var (
	errNotLoggedIn = errors.New("You should login first.")
	timeout        = time.Duration(10 * time.Second)
)

func AuthenticatedSenseClient(auth bool) (*sense.SenseClient, error) {
	return AuthenticatedSenseClientWithTimeout(auth, timeout)
}

func AuthenticatedSenseClientWithTimeout(auth bool, timeout time.Duration) (*sense.SenseClient, error) {
	c := sense.NewClient(nil, timeout)
	if auth {
		home := os.Getenv("HOME")
		settingsPath := filepath.Join(home, SettingsFileName)
		bytes, err := ioutil.ReadFile(settingsPath)
		if err != nil {
			return nil, err
		}

		if len(bytes) == 0 {
			return nil, errNotLoggedIn
		}
		c.AuthToken = string(bytes)
	}
	return c, nil
}

func SenseProtobufClient() (*sense.SenseProtobufClient, error) {
	c := sense.NewProtobufClient(nil, timeout)
	return c, nil
}
