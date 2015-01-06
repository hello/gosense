package sense

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	// "io/ioutil"
	"crypto/hmac"
	"crypto/sha256"
	"net/http"
	"net/url"
)

type AccountService struct {
	client *SenseClient
}

type Gender string

const (
	MALE   Gender = "MALE"
	FEMALE Gender = "FEMALE"
	OTHER  Gender = "OTHER"
)

type Account struct {
	Name         string `json:"name,omitempty"`
	Email        string `json:"email,omitempty"`
	Gender       Gender `json:"gender, omitempty"`
	Height       int32  `json:"height, omitempty"`
	Weight       int32  `json:"weight, omitempty"`
	DOB          string `json:"dob, omitempty"`
	LastModified int64  `json:"last_modified, omitempty"`
}

type Registration struct {
	Name     string  `json:"name,omitempty"`
	Email    string  `json:"email,omitempty"`
	Gender   Gender  `json:"gender, omitempty"`
	Height   int32   `json:"height, omitempty"`
	Weight   int32   `json:"weight, omitempty"`
	TimeZone int32   `json:"tz, omitempty"`
	Password string  `json:"password, omitempty"`
	Lat      float32 `json:"lat, omitempty"`
	Lon      float32 `json:"lat, omitempty"`
}

func NewRegistration(name, email, password string) *Registration {
	return &Registration{Name: name, Email: email, Password: password, Gender: OTHER}
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

func ComputeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (s *AccountService) Register(reg *Registration) (Account, *http.Response, error) {
	// reg := &Registration{
	// 	Name:     "tim",
	// 	Email:    "blah@gmail.com",
	// 	Password: "Oh yeah",
	// 	Gender:   "OTHER",
	// 	TimeZone: -252000,
	// }

	res1B, err := json.Marshal(reg)
	if err != nil {
		return Account{}, nil, err
	}

	body := bytes.NewBuffer(res1B)
	// b, err := ioutil.ReadAll(body)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(b))

	sig := fmt.Sprintf("?sig=%s", url.QueryEscape(ComputeHmac256("timtimtim", "hello")))
	req, err := s.client.NewRequest("POST", "v1/account"+sig, body)
	req.Header.Del("Content-type")
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Date", "Mon, 09 Sep 2011 23:36:00 GMT")

	account := new(Account)
	resp, err := s.client.Do(req, account)
	if err != nil {
		return Account{}, resp, err
	}

	return *account, resp, err
}

func (s *AccountService) Update(account *Account) (Account, *http.Response, error) {
	// reg := &Registration{
	// 	Name:     "tim",
	// 	Email:    "blah@gmail.com",
	// 	Password: "Oh yeah",
	// 	Gender:   "OTHER",
	// 	TimeZone: -252000,
	// }

	account_bytes, err := json.Marshal(account)
	if err != nil {
		return Account{}, nil, err
	}

	body := bytes.NewBuffer(account_bytes)
	req, err := s.client.NewRequest("PUT", "v1/account", body)
	req.Header.Del("Content-type")
	req.Header.Add("Content-type", "application/json")

	updated_account := new(Account)
	resp, err := s.client.Do(req, updated_account)
	if err != nil {
		return Account{}, resp, err
	}

	return *updated_account, resp, err
}
