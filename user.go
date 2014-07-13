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
	usersUrl     = apiUrl + "/users"
	userInfoUrl  = usersUrl + "/show.json"
	userCountUrl = usersUrl + "/counts.json"
)

const (
	getUserInfoModeUid        = iota
	getUserInfoModeScreenName = iota
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
	Status             Status
	Allow_all_comment  bool
	Avatar_large       string
	Verified_reason    string
	Follow_me          bool
	Online_status      int64
	Bi_followers_count int64
}

type Status struct {
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

func (c *Client) getUserInfo(param string, mode int) (User, bool) {
	params := url.Values{"access_token": {c.AccessToken}}
	if mode == getUserInfoModeUid {
		params.Add("uid", param)
	} else {
		params.Add("screen_name", param)
	}
	resp, err := http.Get(encodeParameters(userInfoUrl, params))
	panicError(err)

	if resp.StatusCode >= 400 {
		return User{}, false
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		panicError(err)

		var currentResponse User

		err = json.Unmarshal(body, &currentResponse)
		panicError(err)
		return currentResponse, true
	}
}

func (c *Client) GetUserInfoWithUid(uid string) (User, bool) {
	return c.getUserInfo(uid, getUserInfoModeUid)
}

func (c *Client) GetUserInfoWithScreenName(sName string) (User, bool) {
	return c.getUserInfo(sName, getUserInfoModeScreenName)
}

func (c *Client) GetCurrentUserInfo() (User, bool) {
	return c.GetUserInfoWithUid(c.Uid)
}

func (c *Client) GetUsersFollowersFriendsStatusCounts(uids []string) (map[string]map[string]int64, bool) {
	params := url.Values{"access_token": {c.AccessToken}, "uids": {strings.Join(uids, ",")}}
	resp, err := http.Get(encodeParameters(userCountUrl, params))
	panicError(err)

	if resp.StatusCode >= 400 {
		return make(map[string]map[string]int64), false
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		panicError(err)
		var currentResponse UserCountResponse

		err = json.Unmarshal(body, &currentResponse)
		panicError(err)
		result := make(map[string]map[string]int64)
		for _, currentUC := range currentResponse {
			result[strconv.Itoa(currentUC.Id)] = make(map[string]int64)
			result[strconv.Itoa(currentUC.Id)]["followers_count"] = currentUC.Followers_count
			result[strconv.Itoa(currentUC.Id)]["friends_count"] = currentUC.Friends_count
			result[strconv.Itoa(currentUC.Id)]["statuses_count"] = currentUC.Statuses_count
		}

		return result, true
	}
}
