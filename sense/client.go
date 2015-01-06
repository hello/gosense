// This sense client is heavily inspired by https://github.com/google/go-github
// The Github client is distributed under the BSD-style license found at:
//
// https://github.com/google/go-github/blob/master/LICENSE
//
package sense

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL = "https://dev-api.hello.is/"
	//defaultBaseURL = "http://localhost:9999/"
	userAgent      = "gosense/0.1"
)

// A SenseClient manages communication with the Sense API.
type SenseClient struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests.  Defaults to the public Sense API, but can be
	// set to a domain endpoint to use with hosted cloud platform.  BaseURL should
	// always be specified with a trailing slash.
	BaseURL *url.URL

	// User agent used when communicating with the Sense API.
	UserAgent string
	AuthToken string

	// Services used for talking to different parts of the Sense API.
	Account *AccountService
	Tokens  *TokenService
}

// NewClient returns a new Sense API client. If a nil httpClient is
// provided, http.DefaultClient will be used.  To use API methods which require
// authentication, provide an http.Client that will perform the authentication
// for you (such as that provided by the goauth2 library).
func NewClient(httpClient *http.Client, timeout time.Duration) *SenseClient {
	if httpClient == nil {
		httpClient = &http.Client{
			Transport: &http.Transport{
				Dial: func(network, addr string) (net.Conn, error) {
					conn, err := net.DialTimeout(network, addr, timeout)
					if err != nil {
						return nil, err
					}
					conn.SetDeadline(time.Now().Add(timeout))
					return conn, nil
				},
			},
		}
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &SenseClient{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}
	c.Account = &AccountService{client: c}
	c.Tokens = &TokenService{client: c}
	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the SenseClient.
// Relative URLs should always be specified without a preceding slash.
func (c *SenseClient) NewRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if c.AuthToken != "" {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.AuthToken))
	}

	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

// Do sends an API request and returns the API response.  The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.  If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *SenseClient) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	response := resp

	err = CheckResponse(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}
	return response, err
}

// CheckResponse checks the API response for errors, and returns them if
// present.  A response is considered an error if it has a status code outside
// the 200 range.  API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse.  Any other
// response body will be silently ignored.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}

type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:"message"` // error message
	Code     int32          `json:"code"`    // more detail on individual errors
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %+v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message, r.Code)
}

// type Error struct {
// 	Description string `json:"error_description"` // resource on which the error occurred
// 	ErrorString string `json:"error"`             // field on which the error occurred
// 	Code        int32  `json:"code"`              // validation error code
// }

// func (e *Error) Error() string {
// 	return fmt.Sprintf("HTTP %v error (%v) %v ",
// 		e.Code, e.ErrorString, e.Description)
// }
