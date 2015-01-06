// This sense client is heavily inspired by https://github.com/google/go-github
// The Github client is distributed under the BSD-style license found at:
//
// https://github.com/google/go-github/blob/master/LICENSE
//
package sense

import (
	"bytes"
	proto "code.google.com/p/goprotobuf/proto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/hello/sense/hello"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultServiceBaseURL = "https://dev-in.hello.is/"
	// defaultServiceBaseURL = "http://localhost:5555/"
)

// A SenseClient manages communication with the Sense API.
type SenseProtobufClient struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests.  Defaults to the public Sense API, but can be
	// set to a domain endpoint to use with hosted cloud platform.  BaseURL should
	// always be specified with a trailing slash.
	BaseURL *url.URL

	// User agent used when communicating with the Sense API.
	UserAgent string
	AESKey    string

	Upload *UploadService
}

// NewClient returns a new Sense API client. If a nil httpClient is
// provided, http.DefaultClient will be used.  To use API methods which require
// authentication, provide an http.Client that will perform the authentication
// for you (such as that provided by the goauth2 library).
func NewProtobufClient(httpClient *http.Client, timeout time.Duration) *SenseProtobufClient {
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

	baseURL, _ := url.Parse(defaultServiceBaseURL)

	c := &SenseProtobufClient{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}
	c.Upload = &UploadService{client: c}
	log.Println(c.BaseURL)
	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the SenseClient.
// Relative URLs should always be specified without a preceding slash.
func (c *SenseProtobufClient) NewProtobufRequest(method, urlStr string, body []byte, aesKey string) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	log.Printf("Len body : %d\n", len(body))
	log.Printf("Rel : %s\n", rel)

	u := c.BaseURL.ResolveReference(rel)
	keybytes, _ := hex.DecodeString(aesKey)
	encoded_body, err := encode(body, keybytes)

	buff := &bytes.Buffer{}
	buff.Write(encoded_body)

	req, err := http.NewRequest(method, u.String(), buff)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/octet-stream")
	req.Header.Add("Content-Type", "application/x-protobuf")
	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

// Do sends an API request and returns the API response.  The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.  If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *SenseProtobufClient) Do(req *http.Request, aesKey string) (*hello.SyncResponse, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		log.Println("Do req failed")
		return nil, err
	}

	defer resp.Body.Close()

	response := &hello.SyncResponse{}

	err = CheckProtobufResponse(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		log.Println("Returning error", err)
		return response, err
	}
	buff, err := ioutil.ReadAll(resp.Body)
	keybytes, _ := hex.DecodeString(aesKey)
	pb_resp, err := decode(buff, []byte(keybytes))

	if err != nil {
		log.Println("Decoding error")
		return response, nil
	}

	err = proto.Unmarshal(pb_resp, response)

	return response, err
}

// CheckResponse checks the API response for errors, and returns them if
// present.  A response is considered an error if it has a status code outside
// the 200 range.  API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse.  Any other
// response body will be silently ignored.
func CheckProtobufResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		log.Println("200 status code")
		return nil
	}

	log.Printf("%d\n", r.StatusCode)
	defer r.Body.Close()
	data, _ := ioutil.ReadAll(r.Body)
	return errors.New(string(data))
}

func encode(message []byte, key []byte) ([]byte, error) {
	log.Println("len message", len(message))
	iv := make([]byte, 16)
	for i := 0; i < len(iv); i++ {
		iv[i] = byte(i)
	}

	sha_buf := sha1.Sum(message)

	padded_sha := make([]byte, 32)

	for i, c := range sha_buf {
		padded_sha[i] = c
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, len(padded_sha))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padded_sha)

	// log.Printf("iv: %v\n", iv)
	// log.Printf("sha: %v\n", sha_buf)
	log.Printf("sha: %x\n", sha_buf)
	// log.Printf("Padded sha: %v\n", padded_sha)
	log.Printf("Padded sha: %x\n", padded_sha)
	log.Printf("sig: %x\n", ciphertext)

	// log.Printf("len ciphertext: %v\n", len(ciphertext))
	// log.Printf("len message: %v\n", len(message))
	// log.Printf("len iv: %v\n", len(iv))

	c := [][]byte{message, iv, ciphertext}
	resp := bytes.Join(c, []byte(""))
	// log.Println("Len body in encode", len(resp))
	return resp, nil
}

func decode(body []byte, key []byte) ([]byte, error) {
	log.Println("len message", len(body))
	IV_LENGTH := 16
	SIG_LENGTH := 32

	iv := body[0:IV_LENGTH]
	sig := body[IV_LENGTH : IV_LENGTH+SIG_LENGTH]
	pb := body[IV_LENGTH+SIG_LENGTH:]

	sha_buf := sha1.Sum(pb)

	padded_sha := make([]byte, 32)

	for i, c := range sig {
		padded_sha[i] = c
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, len(padded_sha))

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, padded_sha)

	log.Printf("sha: %x\n", sha_buf)
	log.Printf("ciphertext: %x\n", ciphertext)

	for i, c := range sha_buf {
		if c != ciphertext[i] {
			log.Println("DO NOT MATCH")
			return make([]byte, 0), errors.New("DO NOT MATCH")
		}
	}

	// log.Printf("iv: %v\n", iv)
	// log.Printf("sha: %v\n", sha_buf)
	// log.Printf("sha: %x\n", sha_buf)
	// log.Printf("Padded sha: %v\n", padded_sha)
	// log.Printf("Padded sha: %x\n", padded_sha)
	// log.Printf("ciphertext: %x\n", ciphertext)

	// log.Printf("len ciphertext: %v\n", len(ciphertext))
	// log.Printf("len body: %v\n", len(body))
	// log.Printf("len iv: %v\n", len(iv))

	return pb, nil
}
