package evel

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
)

type VESResponseTestSuite struct {
	suite.Suite
}

func TestVESResponse(t *testing.T) {
	suite.Run(t, new(VESResponseTestSuite))
}

func (s *VESResponseTestSuite) TestNoError() {
	resp := VESResponse{
		CommandList:  []Command{},
		RequestError: map[string]*RequestError{},
	}
	s.False(resp.IsError())
	s.NoError(resp.GetError())

	resp = VESResponse{
		CommandList:  nil,
		RequestError: nil,
	}
	s.False(resp.IsError())
	s.NoError(resp.GetError())
}

func (s *VESResponseTestSuite) TestIsError() {
	resp := VESResponse{
		CommandList:  []Command{},
		RequestError: map[string]*RequestError{"serviceException": {}},
	}
	s.True(resp.IsError())
	s.Error(resp.GetError())
}

func (s *VESResponseTestSuite) TestRequestErrorFormatting() {
	err := RequestError{
		MessageID: "ID1234",
		Text:      "test with var $1 and var $2",
		Variables: []string{"v1", "V2"},
	}
	s.EqualError(&err, "ID1234: test with var v1 and var V2")
}

func (s *VESResponseTestSuite) TestDecodeErrorNoBody() {
	resp := http.Response{
		StatusCode: 199,
	}
	vesresp, err := DecodeVESResponse(&resp)
	s.Nil(vesresp)
	s.Error(err)
}

func (s *VESResponseTestSuite) TestDecodeCommandList() {
	body := bytes.NewBufferString(`
	{
		"commandList": [
			{
				"commandType": "heartbeatIntervalChange",
				"heartbeatInterval": 120
			}
		]
	}
	`)

	resp := http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(body),
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
	}

	vesresp, err := DecodeVESResponse(&resp)
	s.NotNil(vesresp)
	s.NoError(err)
	s.False(vesresp.IsError())
	s.Len(vesresp.CommandList, 1)
	s.Equal(vesresp.CommandList[0].CommandType, CommandHeartbeatIntervalChange)

}

func (s *VESResponseTestSuite) TestDecodeServiceException() {
	body := bytes.NewBufferString(`
	{
		"requestError": {
			"serviceException": {
				"messageId": "SVC12"
			}
		}
	}
	`)

	resp := http.Response{
		StatusCode: 400,
		Body:       ioutil.NopCloser(body),
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
	}

	vesresp, err := DecodeVESResponse(&resp)
	s.Nil(vesresp)
	s.Error(err)
	reqErr, ok := err.(*RequestError)
	s.True(ok)
	s.Equal("SVC12", reqErr.MessageID)

}

func (s *VESResponseTestSuite) TestDecodeInvalidBody() {
	body := bytes.NewBufferString("{")

	resp := http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(body),
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
	}

	vesresp, err := DecodeVESResponse(&resp)
	s.NotNil(vesresp)
	s.NoError(err)
}

type ClientTestSuite struct {
	suite.Suite
}

func TestClient(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

func (s *ClientTestSuite) TestCreateRequest() {
	baseURL, _ := url.Parse("http://foo:bar@1.2.3.4:1234/base")
	client := NewVESClient(*baseURL, nil, nil, 0)
	req, err := client.CreateJSONPostRequest("/api", struct{ Foo string }{"foobar"})
	s.NoError(err)
	s.NotNil(req)
	s.Equal("application/json", req.Header.Get("Content-Type"))
	s.Equal("POST", req.Method)
	user, pass, ok := req.BasicAuth()
	s.True(ok)
	s.Equal("foo", user)
	s.Equal("bar", pass)
	s.Equal("/base/api", req.URL.Path)
	body, err := ioutil.ReadAll(req.Body)
	s.NoError(err)
	s.Equal("{\"Foo\":\"foobar\"}\n", string(body))
}

func (s *ClientTestSuite) TestPostJSON() {
	ncall := 0
	data := map[string]string{"a": "1", "b": "2"}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ncall++
		s.Equal("/base/foobar", req.URL.Path)
		s.Equal("application/json", req.Header.Get("Content-Type"))
		rdata := make(map[string]string)
		err := json.NewDecoder(req.Body).Decode(&rdata)
		s.NoError(err)
		s.Equal(data, rdata)
	}))
	defer srv.Close()
	baseURL, _ := url.Parse(srv.URL + "/base")
	client := NewVESClient(*baseURL, nil, nil, 0)
	resp, err := client.PostJSON("/foobar", data)
	s.NoError(err)
	s.NotNil(resp)
	s.Len(resp.CommandList, 0)
	s.Len(resp.RequestError, 0)
	s.Equal(1, ncall)
}

func (s *ClientTestSuite) TestPostJSONUnreachable() {
	baseURL, _ := url.Parse("http://localhost:1234/base")
	client := NewVESClient(*baseURL, nil, nil, 0)
	resp, err := client.PostJSON("/foobar", nil)
	s.Error(err)
	s.Nil(resp)
}

func (s *ClientTestSuite) TestErrBodyTooLarge() {
	baseURL, _ := url.Parse("http://localhost:1234/base")
	client := NewVESClient(*baseURL, nil, nil, 500)
	resp, err := client.PostJSON("/foobar", make([]int, 1000))
	s.Error(err)
	s.Equal(err, ErrBodyTooLarge)
	s.Nil(resp)
}
