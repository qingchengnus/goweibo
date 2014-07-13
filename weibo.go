package goweibo

import (
	//"net/http"
	"net/url"
)

type Client struct {
	AppKey      string
	AppSecret   string
	CallbackUrl string
	AccessToken string
	Uid         string
}

const (
	baseUrl = "https://api.weibo.com"
	version = 2
	apiUrl  = baseUrl + "/" + string(version)
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
