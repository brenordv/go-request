package models

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brenordv/go-request/internal/core"
	"github.com/brenordv/go-request/internal/parsers"
	"net/http"
	"strings"
)

type HttpConfig struct {
	RequestType         string
	AggressiveMode      bool                   `json:"aggressiveMode"`
	MaxParallelRequests int                    `json:"maxParallelRequests"`
	NumRequests         int                    `json:"numRequests"`
	Url                 string                 `json:"url"`
	QueryString         map[string]interface{} `json:"queryString"`
	Headers             map[string]interface{} `json:"headers"`
	Body                map[string]interface{} `json:"body"`
}

func (h *HttpConfig) GetQueryString() string {
	if h.QueryString == nil {
		return ""
	}
	queryString := parsers.GenericMapToQueryString(h.QueryString)
	return queryString
}

func (h *HttpConfig) MakeRequest() (*http.Request, error) {
	url := h.Url
	queryString := h.GetQueryString()
	if queryString != "" {
		url = fmt.Sprintf("%s?%s", url, queryString)
	}

	var err error
	var req *http.Request

	if h.Body == nil {
		req, err = http.NewRequest(h.RequestType, url, nil)
	} else {
		payloadBuf := new(bytes.Buffer)
		err := json.NewEncoder(payloadBuf).Encode(h.Body)
		if err != nil {
			return nil, err
		}

		req, err = http.NewRequest(h.RequestType, url, payloadBuf)
	}

	if err != nil {
		return nil, err
	}

	if h.Headers != nil {
		for key, value := range h.Headers {
			req.Header.Add(key, parsers.InterfaceToString(value))
		}
	}

	return req, nil
}

func (h *HttpConfig) GetHttpClient() *http.Client {
	client := &http.Client{}

	if strings.HasPrefix(strings.ToLower(h.Url), "https") {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	return client
}

func (h *HttpConfig) Validate() error {
	if h.Url == "" {
		return errors.New("key 'url' cannot be empty or null")
	}

	if h.NumRequests <= 0 {
		return errors.New("key 'numRequests' must be greater than zero")
	}

	if h.MaxParallelRequests <= 0 {
		h.MaxParallelRequests = core.MaxParallelRequests
	}

	if h.RequestType == core.HttpPost && h.Body == nil {
		return errors.New("cannot make a POST request without a body")
	}

	return nil
}
