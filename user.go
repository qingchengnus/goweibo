package goweibo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	usersUrl    = apiUrl + "/users"
	userInfoUrl = usersUrl + "/show.json"
)

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
	Created_at         int64
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

func (c *Client) GetUserInfoWithUid(uid string) (User, bool) {
	params := url.Values{"access_token": {c.AccessToken}, "uid": {uid}}
	resp, err := http.Get(encodeParameters(userInfoUrl, params))
	panicError(err)

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		panicError(err)
		panic(string(body))
		return User{}, false
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		panicError(err)

		var currentResponse User

		err = json.Unmarshal(body, &currentResponse)
		fmt.Println("Body: ", currentResponse)
		panicError(err)
		return currentResponse, true
	}
}
