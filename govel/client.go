/*
	Copyright 2019 Nokia

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package govel

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
	"github.com/nokia/onap-vespa/govel/schema"

	log "github.com/sirupsen/logrus"
)

var (
	// ErrBodyTooLarge is an error returned when JSON body size exceed the configured maximum size (if any)
	ErrBodyTooLarge = errors.New("VES client: Request's body is too large")
)

// RequestError holds details of errors sent back by VES server
type RequestError struct {
	MessageID string   `json:"messageId"`
	Text      string   `json:"text,omitempty"`
	URL       string   `json:"url,omitempty"`
	Variables []string `json:"variables,omitempty"`
}

func (e *RequestError) Error() string {
	text := e.Text
	for i := 0; i < len(e.Variables); i++ {
		text = strings.Replace(text, fmt.Sprintf("$%d", i+1), e.Variables[i], -1)
	}
	return fmt.Sprintf("%s: %s", e.MessageID, text)
}

// VESResponse is the optional response from VES server after an event has been posted
type VESResponse struct {
	CommandList  []Command                `json:"commandList"`
	RequestError map[string]*RequestError `json:"requestError,omitempty"`
}

// IsError return `true` if the response is an error
func (r *VESResponse) IsError() bool {
	return len(r.RequestError) > 0
}

// GetError return the error or nil if not an error
func (r *VESResponse) GetError() error {
	for _, e := range r.RequestError {
		// Return the first error found in map
		return e
	}
	return nil
}

// VESClient is the HTTP client used to talk to VES collector server
type VESClient struct {
	baseURL     url.URL
	client      http.Client
	schema      *schema.JSONSchema
	maxBodySize int
}

// NewVESClient creates a new HTTP client.
// `tlsConfig` is optional, and used only for HTTPS connections.
// `schema` is also optional, and if provided, any outgoing request will have it's JSON payload validated
// `baseURL` should contains user & password
func NewVESClient(baseURL url.URL, tlsConfig *tls.Config, schema *schema.JSONSchema, maxBodySize int) *VESClient {
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       tlsConfig,
		},
	}
	return &VESClient{baseURL: baseURL, client: client, schema: schema, maxBodySize: maxBodySize}
}

// ValidateWithSchema validates the provided data with the client schema.
// If no schema was provided, then validation is silently skipped
func (ves *VESClient) ValidateWithSchema(data interface{}) error {
	if ves.schema != nil {
		log.Debug("Validating request payload with schema before sending it")
		return ves.schema.Validate(data)
	}
	return nil
}

// CreateJSONPostRequest creates an HTTP POST request with `queryPath` added to the client's baseURL.
// If provided, `data` will be serialized into JSON, else if null an empty JSON object will be used
func (ves *VESClient) CreateJSONPostRequest(queryPath string, data interface{}) (*http.Request, error) {
	if data == nil {
		//if data is nil, then replace it by and empty struct
		data = struct{}{}
	}
	if err := ves.ValidateWithSchema(data); err != nil {
		// TODO: Should panic ???
		return nil, err
	}

	url := ves.baseURL // Copy
	url.Path = path.Join(url.Path, queryPath)
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return nil, err
	}

	if ves.maxBodySize > 0 && buf.Len() > ves.maxBodySize {
		log.Warnf("Request body length (%d) exceed the configured maximum (%d)", buf.Len(), ves.maxBodySize)
		return nil, ErrBodyTooLarge
	}

	req, err := http.NewRequest(http.MethodPost, url.String(), &buf)
	if err != nil {
		return nil, err
	}
	if url.User != nil {
		passwd, _ := url.User.Password()
		req.SetBasicAuth(url.User.Username(), passwd)
	}
	req.ContentLength = int64(buf.Len())
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// SendRequest sends the provided HTTP request using the inner HTTP client, and returns
// a VESResponse, or an error
func (ves *VESClient) SendRequest(req *http.Request) (*VESResponse, error) {
	u := *req.URL
	if u.User != nil {
		// Remove password before printing debug message
		u.User = url.User(u.User.Username())
	}
	log.Debug("Send POST to ", u.String())
	//TODO: Add context
	resp, err := ves.client.Do(req)
	if err != nil {
		return nil, err
	}
	return DecodeVESResponse(resp)
}

// PostJSON sends an HTTP POST request with `queryPath` added to the client's baseURL.
// If provided, `data` is serialized into JSON, else if null an empty JSON object is be used
func (ves *VESClient) PostJSON(queryPath string, data interface{}) (*VESResponse, error) {
	req, err := ves.CreateJSONPostRequest(queryPath, data)
	if err != nil {
		return nil, err
	}
	return ves.SendRequest(req)
}

// DecodeVESResponse transform and http API response into a VESResponse
// or return an error if the response is an error
func DecodeVESResponse(resp *http.Response) (*VESResponse, error) {
	if resp.Body != nil {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Error(err.Error())
			}
		}()
	}
	vesResp := new(VESResponse)
	if resp.Header.Get("Content-Type") == "application/json" {
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(vesResp); err != nil {
			log.Warn("Could not decode JSON response: ", err.Error())
		}
		log.Debugf("Got response %+v", vesResp)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if vesResp.IsError() {
			return nil, vesResp.GetError()
		}
		return nil, fmt.Errorf("HTTP request failed (status %d)", resp.StatusCode)
	}
	return vesResp, nil
}
