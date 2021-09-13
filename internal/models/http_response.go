package models

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type HttpResponse struct {
	SessionId             string
	ReqNumber             int
	Url                   string
	ParsedUrl             string
	StatusCode            int
	Status                string
	ReqHeaders            map[string][]string
	ResHeaders            map[string][]string
	Body                  string
	ReadBodyError         string
	Error                 string
	IsAuthenticationError bool
	IsAuthorizationError  bool
	IsInternalServerError bool
	Ok                    bool
}

func getHttpResponseBody(res *http.Response) (string, error) {
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), err
}

func headersToString(headers map[string][]string) string {
	var headerLines []string
	for name, values := range headers {
		for _, value := range values {
			headerLines = append(headerLines, fmt.Sprintf("%s: %v", name, value))
		}
	}
	return strings.Join(headerLines, "\n")
}

func NewHttpResponse(url string, sessionId string, reqNum int, res *http.Response, err error) *HttpResponse {
	var response HttpResponse
	if res != nil {
		body, bErr := getHttpResponseBody(res)
		response = HttpResponse{
			ParsedUrl:             res.Request.URL.String(),
			StatusCode:            res.StatusCode,
			Status:                res.Status,
			ReqHeaders:            res.Request.Header,
			ResHeaders:            res.Header,
			Body:                  body,
			ReadBodyError:         fmt.Sprintf("%v", bErr),
			IsAuthenticationError: res.StatusCode == 401,
			IsAuthorizationError:  res.StatusCode == 403,
			IsInternalServerError: res.StatusCode >= 500 && res.StatusCode <= 599,
		}
	} else {
		response = HttpResponse{}
	}

	response.Url = url
	response.SessionId = sessionId
	response.ReqNumber = reqNum
	response.Error = fmt.Sprintf("%v", err)
	response.Ok = err == nil && res.StatusCode >= 200 && res.StatusCode <= 299
	return &response
}

func DeserializeHttpResponse(data []byte) (*HttpResponse, error) {
	var response HttpResponse
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&response)
	return &response, err
}

func (r *HttpResponse) Serialize() ([]byte, []byte, error) {
	var value bytes.Buffer
	if r.Url == "" {
		return nil, nil, errors.New("url cannot be empty")
	}

	key := []byte(fmt.Sprintf("%s|%s", r.SessionId, r.Url))
	encoder := gob.NewEncoder(&value)
	err := encoder.Encode(r)
	return key, value.Bytes(), err
}

func (r *HttpResponse) String() string {
	return fmt.Sprintf(`-|session: %s|--------------------------------------------------------------------------
[Req #%07d] Url: %s
StatusCode: %d | %s
Success: %v
---Request Headers:
%s

---Request Headers:
%s

---Body:
%s

---Error:
%v
	`,
		r.SessionId,
		r.ReqNumber, r.Url,
		r.StatusCode, r.Status,
		r.Ok,
		headersToString(r.ReqHeaders),
		headersToString(r.ResHeaders),
		r.Body,
		r.Error)
}
