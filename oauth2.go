package goweibo

import (
	"encoding/json"
	"errors"
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

func (c Client) getAuthorizationParameters() (url.Values, error) {
	if c.AppKey == "" || c.CallbackUrl == "" {
		return nil, errors.New("AppKey or CallbackUrl not set!")
	} else {
		return url.Values{"client_id": {c.AppKey}, "redirect_uri": {c.CallbackUrl}}, nil
	}
}

func (c *Client) GetAuthorizationUrl() (string, error) {
	params, err := (*c).getAuthorizationParameters()
	if err != nil {
		return "", err
	}
	return encodeParameters(baseUrl+authorizationUrl, params), nil
}

func (c Client) getAccessTokenParameters(code string) url.Values {
	return url.Values{"client_id": {c.AppKey}, "client_secret": {c.AppSecret}, "grant_type": {"authorization_code"}, "code": {code}, "redirect_uri": {c.CallbackUrl}}

}

func (c *Client) RequestAccessToken(code string) (string, string, string, string, error) {
	resp, httpError := http.PostForm(baseUrl+accessTokenUrl, (*c).getAccessTokenParameters(code))
	if httpError != nil {
		return "", "", "", "", httpError
	}

	defer resp.Body.Close()
	body, readError := ioutil.ReadAll(resp.Body)
	if readError != nil {
		return "", "", "", "", readError
	}
	if resp.StatusCode >= 400 {
		var currentResponse ErrorResponse
		jsonError := json.Unmarshal(body, &currentResponse)
		if jsonError != nil {
			return "", "", "", "", jsonError
		}

		return "", "", "", "", currentResponse.parseErrorResponse()
	} else {

		var currentResponse accessTokenResponse

		jsonError := json.Unmarshal(body, &currentResponse)
		if jsonError != nil {
			return "", "", "", "", jsonError
		}
		c.AccessToken = currentResponse.Access_token
		c.Uid = currentResponse.Uid
		return currentResponse.Access_token, string(currentResponse.Expires_in), currentResponse.Remind_in, currentResponse.Uid, nil
	}
}
