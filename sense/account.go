package sense

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AccountService struct {
	client *SenseClient
}

type Account struct {
	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
	Gender string `json:"gender, omitempty"`
	Height int32  `json:"height, omitempty"`
	Weight int32  `json:"weight, omitempty"`
	DOB    int64  `json:"dob, omitempty"`
}

type Registration struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Gender   string `json:"gender, omitempty"`
	Height   int32  `json:"height, omitempty"`
	Weight   int32  `json:"weight, omitempty"`
	TimeZone int32  `json:"tz, omitempty"`
	Password string `json:"password, omitempty"`
}

func (a *Account) String() string {
	return fmt.Sprintf("Account: %s [%s]", a.Name, a.Email)
}

func (s *AccountService) Me() (Account, *http.Response, error) {

	req, err := s.client.NewRequest("GET", "v1/account", nil)
	if err != nil {
		return Account{}, nil, err
	}

	account := new(Account)
	resp, err := s.client.Do(req, account)
	if err != nil {
		return Account{}, resp, err
	}

	return *account, resp, err
}

func (s *AccountService) Register() (Account, *http.Response, error) {
	reg := &Registration{
		Name:     "tim",
		Email:    "blah@gmail.com",
		Password: "Oh yeah",
		Gender:   "OTHER",
		TimeZone: -252000,
	}

	res1B, err := json.Marshal(reg)
	if err != nil {
		return Account{}, nil, err
	}

	body := bytes.NewBuffer(res1B)
	req, err := s.client.NewRequest("POST", "v1/account", body)
	req.Header.Del("Content-type")
	req.Header.Add("Content-type", "application/json")

	account := new(Account)
	resp, err := s.client.Do(req, account)
	if err != nil {
		return Account{}, resp, err
	}

	return *account, resp, err
}
