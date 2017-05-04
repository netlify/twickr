package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"strings"

	"github.com/darkhelmet/twitterstream"
)

var config *Config

// Config holds all the configuration data
type Config struct {
	ConsumerKey      string `json:"consumerKey"`
	ConsumerSecret   string `json:"consumerSecret"`
	AccessToken      string `json:"accessToken"`
	AccessSecret     string `json:"accessSecret"`
	SlackWebhook     string `json:"slackWebhook"`
	SlackChannel     string `json:"slackChannel"`
	Keywords         string `json:"keywords"`
	KeywordsIgnored  string `json:"keywordsIgnored"`
}

type attachment map[string]string

// SlackMsg represents a message sent to Slack over a webhook
type SlackMsg struct {
	User        string       `json:"username"`
	Text        string       `json:"text"`
	Channel     *string      `json:"channel"`
	Attachments []attachment `json:"attachments"`
}

func init() {
	config = &Config{}
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Error reading configuration: %v", err)
	}
	err = json.Unmarshal(data, config)
	if err != nil {
		log.Fatalf("Error parsing configuration: %v", err)
	}
}

func decode(conn *twitterstream.Connection) {
	for {
		if tweet, err := conn.Next(); err == nil {
			if !testIgnoredKeywords(tweet.Text) {
				continue
			}
			msg := createMessage(tweet)
			b, err := json.Marshal(msg)
			if err != nil {
				log.Printf("Error encoding json for tweet: %v %v\n", tweet, err)
			} else {
				buf := bytes.NewReader(b)
				http.Post(config.SlackWebhook, "application/json", buf)
			}

			log.Printf("%s said: %s\n", tweet.User.ScreenName, tweet.Text)
		} else {
			log.Printf("Failed decoding tweet: %s", err)
			return
		}
	}
}

func testIgnoredKeywords(tweetText string) bool {
	if config.KeywordsIgnored != "" {
		for _, keyword := range strings.Split(config.KeywordsIgnored, ",") {
			if strings.Contains(tweetText, keyword) {
				log.Printf("Tweet contained ignored keyword: %s", keyword)
				return false
			}
		}
	}
	return true
}

func createMessage(tweet *twitterstream.Tweet) *SlackMsg {
	text := "https://twitter.com/" + tweet.User.IdStr + "/status/" + tweet.IdString

	att := attachment{
		"color":       "grey",
		"author_icon": tweet.User.ProfileImageUrlHttps,
		"author_name": fmt.Sprintf("%s @%s", tweet.User.Name, tweet.User.ScreenName),
		"author_link": fmt.Sprintf("https://twitter.com/%s", tweet.User.ScreenName),
		"text":        tweet.Text,
	}

	if len(tweet.Entities.Media) > 0 {
		att["image_url"] = tweet.Entities.Media[0].SecureMediaUrl
	}

	return &SlackMsg{
		User:        tweet.User.ScreenName,
		Text:        text,
		Attachments: []attachment{att},
		Channel:     &config.SlackChannel,
	}
}

func main() {
	client := twitterstream.NewClient(
		config.ConsumerKey,
		config.ConsumerSecret,
		config.AccessToken,
		config.AccessSecret,
	)
	for {
		conn, err := client.Track(config.Keywords)
		if err != nil {
			log.Printf("Tracking failed, sleeping for 1 minute: %v", err)
			time.Sleep(1 * time.Minute)
			continue
		}
		decode(conn)
	}
}
