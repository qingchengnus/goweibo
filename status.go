package goweibo

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	statusesUrl       = apiUrl + "/statuses"
	publicTimelineUrl = statusesUrl + "/public_timeline.json"
)

type Status struct {
	StatusSimple
	User User
}

type StatusesResponse struct {
	Statuses []Status
}

func (c *Client) GetPublicTimeline(numberOfStatuses int) ([]Status, error) {
	if numberOfStatuses > 200 {
		return nil, errors.New("Cannot fetch more than 200 statuses in one request.")
	}
	params := url.Values{"access_token": {c.AccessToken}, "count": {strconv.Itoa(numberOfStatuses)}}
	resp, httpError := http.Get(encodeParameters(publicTimelineUrl, params))
	if httpError != nil {
		return nil, httpError
	}
	defer resp.Body.Close()
	body, readError := ioutil.ReadAll(resp.Body)
	if readError != nil {
		return nil, readError
	}

	if resp.StatusCode >= 400 {
		var errorResponse ErrorResponse
		jsonError := json.Unmarshal(body, &errorResponse)
		if jsonError != nil {
			return nil, jsonError
		}
		return nil, errorResponse.parseErrorResponse()

	} else {

		var currentResponse StatusesResponse

		jsonError := json.Unmarshal(body, &currentResponse)
		if jsonError != nil {
			return nil, jsonError
		}
		return currentResponse.Statuses, nil
	}
}

func (c *Client) GetPublicTimelineDefaultCount() ([]Status, error) {
	return c.GetPublicTimeline(20)
}
