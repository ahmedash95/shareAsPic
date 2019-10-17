package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/joho/godotenv"
)

var TWITTER_API_KEY = ""
var TWITTER_API_SECRET = ""
var TWITTER_ACCESS_TOKEN_KEY = ""
var TWITTER_ACCESS_TOKEN_SECRET = ""
var PIC_STORAGE_PATH = ""
var PIC_STORAGE_URL = ""

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	TWITTER_API_KEY = os.Getenv("TWITTER_API_KEY")
	TWITTER_API_SECRET = os.Getenv("TWITTER_API_SECRET")
	TWITTER_ACCESS_TOKEN_KEY = os.Getenv("TWITTER_ACCESS_TOKEN_KEY")
	TWITTER_ACCESS_TOKEN_SECRET = os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	PIC_STORAGE_PATH = os.Getenv("PIC_STORAGE_PATH")
	PIC_STORAGE_URL = os.Getenv("PIC_STORAGE_URL")

	// init chrome client for screenshots
	initChromedpClient()
	// initialize twitter configuration for streaming
	initTwitterClient()
	// initialize redis client
	initRedisClient()
	// initialize logger
	initLogger()

	Logger.Info("App starts: Waiting for tweets")
	fmt.Println("Listening for tweets")

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		processTweet(*tweet)
	}

	params := &twitter.StreamFilterParams{
		Track:         []string{"@ShareAsPic"},
		StallWarnings: twitter.Bool(true),
	}

	stream, err := client.Streams.Filter(params)
	for message := range stream.Messages {
		demux.Handle(message)
	}

}
