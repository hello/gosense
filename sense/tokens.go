package sense

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type TokenService struct {
	client *SenseClient
}

type AccessToken struct {
	Value     string `json:"access_token,omitempty"`
	Type      string `json:"token_type,omitempty"`
	ExpiresIn uint32 `json:"expires_in,omitempty"`
}

func (a AccessToken) String() string {
	return fmt.Sprintf("AccessToken{%s, %s, %d}", a.Value, a.Type, a.ExpiresIn)
}

func (s *TokenService) Login(username, password string) (AccessToken, *http.Response, error) {
	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	data.Set("grant_type", "password")
	data.Set("client_id", "iphone_pill")
	data.Set("client_secret", "client_secret_here")

	body := bytes.NewBufferString(data.Encode())

	req, err := s.client.NewRequest("POST", "v1/oauth2/token", body)

	token := new(AccessToken)
	resp, err := s.client.Do(req, token)
	log.Println(resp)
	if err != nil {
		return AccessToken{}, resp, err
	}

	return *token, resp, err
}

func (s *TokenService) Delete(token, username, password string) (*http.Response, error) {
	path := fmt.Sprintf("v1/access_tokens/%s", token)
	req, err := s.client.NewRequest("DELETE", path, nil)

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(username, password)

	resp, err := s.client.Do(req, nil)
	return resp, err
}
