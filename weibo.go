package goweibo

import (
	//"net/http"
	"net/url"
)

type Client struct {
	AppKey      string
	AppSecret   string
	CallbackUrl string
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
