package main

import (
	"log"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/joho/godotenv"
)

var twitterAPIKey = ""
var twitterAPISECRET = ""
var twitterAccessTokenKey = ""
var twitterAccessTokenSecret = ""
var picStoragePath = ""

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	twitterAPIKey = os.Getenv("TWITTER_API_KEY")
	twitterAPISECRET = os.Getenv("TWITTER_API_SECRET")
	twitterAccessTokenKey = os.Getenv("TWITTER_ACCESS_TOKEN_KEY")
	twitterAccessTokenSecret = os.Getenv("TWITTER_ACCESS_TOKEN_SECRET")
	picStoragePath = os.Getenv("PIC_STORAGE_PATH")

	// initialize logger
	initLogger()
	// initialize twitter configuration for streaming
	logAndPrint("Init Twitter client")
	initTwitterClient()
	// initialize redis client
	logAndPrint("Init redis")
	initRedisClient()

	logAndPrint("App starts: Waiting for tweets")

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
