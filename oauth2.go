package goweibo

import (
	"encoding/json"
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
		return url.Values{"client_id": {c.AppKey}, "redirect_uri": {c.CallbackUrl}}
	}
}

func (c *Client) GetAuthorizationUrl() string {
	return encodeParameters(baseUrl+authorizationUrl, (*c).getAuthorizationParameters())
}

func (c Client) getAccessTokenParameters(code string) url.Values {
	return url.Values{"client_id": {c.AppKey}, "client_secret": {c.AppSecret}, "grant_type": {"authorization_code"}, "code": {code}, "redirect_uri": {c.CallbackUrl}}

}

func (c *Client) RequestAccessToken(code string) (string, string, string, string, bool) {
	resp, err := http.PostForm(baseUrl+accessTokenUrl, (*c).getAccessTokenParameters(code))
	panicError(err)

	if resp.StatusCode >= 400 {
		return "", "", "", "", false
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		panicError(err)

		var currentResponse accessTokenResponse

		err = json.Unmarshal(body, &currentResponse)
		panicError(err)
		c.SetAccessToken(currentResponse.Access_token)
		c.Uid = currentResponse.Uid
		return currentResponse.Access_token, string(currentResponse.Expires_in), currentResponse.Remind_in, currentResponse.Uid, true
	}
}

func (c *Client) SetAccessToken(accessToken string) {
	(*c).AccessToken = accessToken
}
