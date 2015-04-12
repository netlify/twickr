package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/darkhelmet/twitterstream"
)

var config *Config

// Config holds all the configuration data
type Config struct {
	ConsumerKey    string `json:"consumerKey"`
	ConsumerSecret string `json:"consumerSecret"`
	AccessToken    string `json:"accessToken"`
	AccessSecret   string `json:"accessSecret"`
	SlackWebhook   string `json:"slackWebhook"`
	Keywords       string `json:"keywords"`
}

// SlackMsg represents a message sent to Slack over a webhook
type SlackMsg struct {
	User string `json:"username"`
	Text string `json:"text"`
}

func init() {
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
			msg := &SlackMsg{
				User: tweet.User.ScreenName,
				Text: "https://twitter.com/" + tweet.User.IdStr + "/status/" + tweet.IdString,
			}
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
			log.Println("Tracking failed, sleeping for 1 minute")
			time.Sleep(1 * time.Minute)
			continue
		}
		decode(conn)
	}
}
