package goweibo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	usersUrl          = apiUrl + "/users"
	userInfoUrl       = usersUrl + "/show.json"
	userCountUrl      = usersUrl + "/counts.json"
	userInfoDomainUrl = usersUrl + "/domain_show.json"
)

const (
	getUserInfoModeUid        = iota
	getUserInfoModeScreenName = iota
	getUserInfoModeDomain     = iota
)

type UserCountResponse []UserCount

type UserCount struct {
	Id              int64
	Followers_count int64
	Friends_count   int64
	Statuses_count  int64
}

type User struct {
	Id                 int64
	Screen_name        string
	Name               string
	Province           string
	City               string
	Location           string
	Description        string
	Url                string
	Profile_image_url  string
	Domain             string
	Gender             string
	Followers_count    int64
	Friends_count      int64
	Statuses_count     int64
	Favourites_count   int64
	Created_at         string
	Following          bool
	Allow_all_act_msg  bool
	Geo_enabled        bool
	Verified           bool
	Status             StatusSimple
	Allow_all_comment  bool
	Avatar_large       string
	Verified_reason    string
	Follow_me          bool
	Online_status      int64
	Bi_followers_count int64
}

type StatusSimple struct {
	Created_at              string
	Id                      int64
	Text                    string
	Source                  string
	Favorited               bool
	Truncated               bool
	In_reply_to_status_id   string
	In_reply_to_user_id     string
	In_reply_to_screen_name string
	Geo                     Geo
	Mid                     string
	Reposts_count           int64
	Comments_count          int64
}

type Geo struct {
	Longitude     string
	Latitude      string
	City          string
	Province      string
	City_name     string
	Province_name string
	Address       string
	Pinyin        string
	More          string
}

func (c *Client) getUserInfo(param string, mode int) (User, error) {
	params := url.Values{"access_token": {c.AccessToken}}
	var requestUrl string
	if mode == getUserInfoModeUid {
		params.Add("uid", param)
		requestUrl = userInfoUrl
	} else if mode == getUserInfoModeScreenName {
		params.Add("screen_name", param)
		requestUrl = userInfoUrl
	} else {
		params.Add("domain", param)
		requestUrl = userInfoDomainUrl
	}

	resp, httpError := http.Get(encodeParameters(requestUrl, params))
	if httpError != nil {
		return User{}, httpError
	}
	defer resp.Body.Close()
	body, readError := ioutil.ReadAll(resp.Body)
	if readError != nil {
		return User{}, readError
	}
	if resp.StatusCode >= 400 {
		var currentResponse ErrorResponse

		jsonError := json.Unmarshal(body, &currentResponse)
		if jsonError != nil {
			return User{}, jsonError
		}
		return User{}, currentResponse.parseErrorResponse()
	} else {

		var currentResponse User

		jsonError := json.Unmarshal(body, &currentResponse)
		if jsonError != nil {
			return User{}, jsonError
		}
		return currentResponse, nil
	}
}

func (c *Client) GetUserInfoWithUid(uid string) (User, error) {
	return c.getUserInfo(uid, getUserInfoModeUid)
}

func (c *Client) GetUserInfoWithScreenName(sName string) (User, error) {
	return c.getUserInfo(sName, getUserInfoModeScreenName)
}

func (c *Client) GetUserInfoWithDomain(domain string) (User, error) {
	return c.getUserInfo(domain, getUserInfoModeDomain)
}

func (c *Client) GetCurrentUserInfo() (User, error) {
	return c.GetUserInfoWithUid(c.Uid)
}

func (c *Client) GetUsersFollowersFriendsStatusCounts(uids []string) (map[string]map[string]int64, error) {
	params := url.Values{"access_token": {c.AccessToken}, "uids": {strings.Join(uids, ",")}}
	resp, httpError := http.Get(encodeParameters(userCountUrl, params))
	if httpError != nil {
		return nil, httpError
	}

	defer resp.Body.Close()
	body, readError := ioutil.ReadAll(resp.Body)

	if readError != nil {
		return nil, readError
	}

	if resp.StatusCode >= 400 {
		var currentResponse ErrorResponse

		jsonError := json.Unmarshal(body, &currentResponse)
		if jsonError != nil {
			return nil, jsonError
		}
		return nil, currentResponse.parseErrorResponse()
	} else {
		var currentResponse UserCountResponse

		jsonError := json.Unmarshal(body, &currentResponse)
		if jsonError != nil {
			return nil, jsonError
		}
		result := make(map[string]map[string]int64)
		for _, currentUC := range currentResponse {
			key := strconv.FormatInt(currentUC.Id, 10)
			result[key] = make(map[string]int64)
			result[key]["followers_count"] = currentUC.Followers_count
			result[key]["friends_count"] = currentUC.Friends_count
			result[key]["statuses_count"] = currentUC.Statuses_count
		}

		return result, nil
	}
}
