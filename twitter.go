package main

// Handles optional posting to twitter of everything your markov-bot says
//
// YOU WILL WANT TO BE VERY CAREFUL ABOUT THIS
// It will/may contain internal links, names, other identifying info from what people
// tend to think is private data
//
// Requires a consumer key and secret and then an associated access token and secret
// You can get these from https://dev.twitter.com/apps/new

import (
	"encoding/json"
	"github.com/mrjones/oauth"
	"io/ioutil"
	"net/http"
)

type Twitter struct {
	consumerKey    string
	consumerSecret string

	accessToken oauth.AccessToken
	client      oauth.Consumer
}

type User struct {
	ID         uint64 `json:"id"`
	ScreenName string `json:"screen_name"`
}

func NewTwitter(consumerKey string, consumerSecret string, accessToken string, accessTokenSecret string) *Twitter {
	c := oauth.NewConsumer(
		consumerKey,
		consumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		})

	return &Twitter{
		consumerKey:    consumerKey,
		consumerSecret: consumerSecret,

		accessToken: oauth.AccessToken{
			Token:  accessToken,
			Secret: accessTokenSecret,
		},

		client: *c,
	}
}

func (t *Twitter) GetMe() (*User, error) {
	resp, err := t.client.Get(
		"https://api.twitter.com/1.1/account/verify_credentials.json",
		map[string]string{},
		&t.accessToken)

	if err != nil {
		return nil, err
	}

	// Make sure we close the body stream no matter what
	defer resp.Body.Close()

	// Read body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// parse into json
	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (t *Twitter) Post(status string) (*http.Response, error) {
	resp, err := t.client.Post(
		"https://api.twitter.com/1.1/statuses/update.json",
		map[string]string{
			"status": status,
		},
		&t.accessToken)

	if err != nil {
		return resp, err
	}

	// Make sure we close the body stream no matter what
	defer resp.Body.Close()

	return resp, nil
}
