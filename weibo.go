package goweibo

import (
	//"net/http"
	"errors"
	"net/url"
)

type Client struct {
	AppKey      string
	AppSecret   string
	CallbackUrl string
	AccessToken string
	Uid         string
}

type ErrorResponse struct {
	Request    string
	Error_code string
	Error      string
}

const (
	baseUrl = "https://api.weibo.com"
	version = "2"
	apiUrl  = baseUrl + "/" + version
)

func encodeParameters(url string, params url.Values) string {
	return url + "?" + params.Encode()
}

func panicError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func NewClient(appKey, appSecret, callbackUrl string) *Client {
	c := new(Client)
	c.AppKey = appKey
	c.AppSecret = appSecret
	c.CallbackUrl = callbackUrl
	return c
}

func (errorResponse ErrorResponse) parseErrorResponse() error {
	return errors.New("Failed to request: " + errorResponse.Request + " . Error code: " + errorResponse.Error_code + " (Refer to http://open.weibo.com/wiki/Error_code). Reason: " + errorResponse.Error + ".")
}
