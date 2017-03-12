package manager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"

	"github.com/aggrolite/rereddit/config"
	"github.com/jroimartin/gocui"
	"github.com/jzelinskie/geddit"
)

type ReRedditManager struct {
	API   *geddit.OAuthSession
	Views []*gocui.View
}

func NewReRedditManager() (*ReRedditManager, error) {
	var err error

	m := new(ReRedditManager)

	// Create reddit API object required for other views.
	m.API, err = newRedditAPI()
	if err != nil {
		return m, err
	}
	return m, nil
}

func newRedditAPI() (*geddit.OAuthSession, error) {
	// For now, user credentials are used for authentication.
	// TODO support application-side, implicit auth?
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	configPath := fmt.Sprintf("%s/.rereddit.json", usr.HomeDir)
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read config file containing credentials. Does %s exist? Error: %s\n", configPath, err)
	}

	c := new(config.Config)
	if err := json.Unmarshal(file, c); err != nil {
		return nil, fmt.Errorf("Error parsing JSON config: %s\n", err)
	}

	o, err := geddit.NewOAuthSession(
		c.ClientID,
		c.ClientSecret,
		"rereddit: a terminal reddit client by u/aggrolite. See source code @ github.com/aggrolite/rereddit",
		"http://redirect.url",
	)
	if err != nil {
		return o, err
	}
	if err := o.LoginAuth(c.User, c.Password); err != nil {
		return o, fmt.Errorf("Failed to login: %s\n", err)
	}
	return o, nil
}
