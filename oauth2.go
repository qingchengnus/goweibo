package goweibo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type accessTokenResponse struct {
	Access_token string
	Expires_in   int64
	Remind_in    string
	Uid          string
}

const (
	authorizationUrl = "/oauth2/authorize"
	accessTokenUrl   = "/oauth2/access_token"
)

const (
	AuthorizationPageDisplayModeDefault    = iota
	AuthorizationPageDisplayModeMobile     = iota
	AuthorizationPageDisplayModeWap        = iota
	AuthorizationPageDisplayModeClient     = iota
	AuthorizationPageDisplayModeApponweibo = iota
)

const (
	AuthorizationLanguageChinese = iota
	AuthorizationLanguageEnglish = iota
)

func (c Client) getAuthorizationParameters() url.Values {
	if c.AppKey == "" || c.CallbackUrl == "" {
		panic("AppKey or CallbackUrl not set!")
	} else {
		//return "?client_id=" + c.AppKey + "&redirect_uri=" + c.CallbackUrl
		return url.Values{"client_id": {c.AppKey}, "redirect_uri": {c.CallbackUrl}}
	}
}

func (c *Client) GetAuthorizationUrl() string {
	//return baseUrl + authorizationUrl + c.getAuthorizationParameters()
	return encodeParameters(baseUrl+authorizationUrl, (*c).getAuthorizationParameters())
}

func (c Client) getAccessTokenParameters(code string) url.Values {
	return url.Values{"client_id": {c.AppKey}, "client_secret": {c.AppSecret}, "grant_type": {"authorization_code"}, "code": {code}, "redirect_uri": {c.CallbackUrl}}

}

func (c *Client) RequestAccessToken(code string) (bool, string, string, string, string) {
	resp, err := http.PostForm(baseUrl+accessTokenUrl, (*c).getAccessTokenParameters(code))
	panicError(err)

	if resp.StatusCode >= 400 {
		return false, "", "", "", ""
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		panicError(err)

		var currentResponse accessTokenResponse

		err = json.Unmarshal(body, &currentResponse)
		fmt.Println("Body: ", currentResponse)
		panicError(err)
		return true, currentResponse.Access_token, string(currentResponse.Expires_in), currentResponse.Remind_in, currentResponse.Uid
	}
}

func (c *Client) SetAccessToken(accessToken string) {
	(*c).AccessToken = accessToken
}
