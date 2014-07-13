goweibo
=======
goweibo is a stand-alone golang library which provides

* OAuth2.0 authentication
* API calling wrapper(Currently only user API supported)

for Sina Weibo 新浪微博.

https://github.com/qingchengnus/goweibo

Developed by Qing Cheng. Contact me at [qingchengnus@gmail.com](qingchengnus@gmail.com).

## Basic Usage

Instantiate client

```go
client := goweibo.NewClient(appKey, appSecret, callbackUrl)
```

## OAuth2.0


Get authorization url
```go
authorizationUrl := client.GetAuthorizationUrl()
```

Request access token
```go
accessToken, expiresIn, remindIn, uid, ok := client.RequestAccessToken(code)
```
where you got the code from callback after user logged in. After this call, you can access current user's uid and access token by calling:
```go
accessToken := client.AccessToken
uid := client.Uid
```


## User API


Get current user info
```go
currentUser := client.GetCurrentUserInfo()
```

Get user info with uid, name or domain
```go
userInfo := client.GetUserInfoWithUid(uid)
userInfo := client.GetUserInfoWithScreenName(name)
userInfo := client.GetUserInfoWithDomain(domain)
name := userInfo.ScreenName
```

Get several users' followers count, friends count and statuses count
```go
usersCount := client.GetUsersFollowersFriendsStatusCounts(uids)
user1FriendsCount := usersCount[uid1]["friends_count"]
```

## wiki

For more information, please refer to standard go 
[docs](https://github.com/michaelliao/goweibo/wiki).
