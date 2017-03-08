package internal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Request struct {
	Url      string            `json:"url"`
	Method   string            `json:"method"`
	Header   map[string]string `json:"header"`
	Playload string            `json:"playload"`
}

func NewRequest(urlStr string, method string) *Request {
	return &Request{
		Url:    urlStr,
		Method: method,
		Header: make(map[string]string),
	}
}

func NewRequestJson(bs []byte) (*Request, error) {
	var r *Request
	err := json.Unmarshal(bs, &r)
	return r, err
}

func (r *Request) BasicAuth(name, psw string) {
	b := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", name, psw)))
	r.Header["Authorization"] = fmt.Sprintf("Basic %s", b)
}

func (r *Request) AsHttpRequest() (*http.Request, error) {
	req, err := http.NewRequest(r.Method, r.Url, strings.NewReader(r.Playload))
	if err != nil {
		return nil, err
	}
	for k, v := range r.Header {
		req.Header.Set(k, v)
	}
	return req, nil
}
