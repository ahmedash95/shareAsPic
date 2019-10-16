package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/go-redis/redis/v7"
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

	for {
		// start streaming
		search, _, err := client.Search.Tweets(&twitter.SearchTweetParams{
			Query: "@ShareAsPic",
		})

		if err != nil {
			Logger.Error(fmt.Sprintf("Faild to fetch latest tweets, %s", err.Error()))
		}

		scanTweets(search)

		// check tweets every minute for 100K/24h request api limit
		time.Sleep(time.Minute)
	}

}

var client *twitter.Client

func initTwitterClient() {
	config := oauth1.NewConfig(TWITTER_API_KEY, TWITTER_API_SECRET)
	token := oauth1.NewToken(TWITTER_ACCESS_TOKEN_KEY, TWITTER_ACCESS_TOKEN_SECRET)
	httpClient := config.Client(oauth1.NoContext, token)
	client = twitter.NewClient(httpClient)
}

var redisClient *redis.Client

func initRedisClient() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		Logger.Error(fmt.Sprintf("Faild to connect to redis, %s", err.Error()))
		log.Fatal(err)
	}
}

func scanTweets(search *twitter.Search) {
	for _, tweet := range search.Statuses {
		if TweetProcessedBefore(tweet) {
			continue
		}
		// let's make sure it has "share this" in the string
		if !strings.Contains(strings.ToLower(tweet.Text), "@shareaspic") {
			continue
		}
		if !strings.Contains(strings.ToLower(tweet.Text), "share this") {
			replyWithIDoNotUnderstand(tweet)
			continue
		}

		makeTweetPicAndShare(tweet)
	}
}

const ProcessedTweets = "processed_tweets"

func TweetProcessedBefore(tweet twitter.Tweet) bool {
	result, _ := redisClient.SAdd(ProcessedTweets, tweet.ID).Result()
	var processed bool
	processed = result == 0
	if !processed {
		log.Printf("new tweet to process: %s\n", tweet.IDStr)
	}

	return processed
}

func replyWithIDoNotUnderstand(tweet twitter.Tweet) {
	log.Printf("replyWithIDoNotUnderstand: %s\n", tweet.IDStr)
	statusUpdate := &twitter.StatusUpdateParams{
		Status:             "",
		InReplyToStatusID:  tweet.ID,
		PossiblySensitive:  nil,
		Lat:                nil,
		Long:               nil,
		PlaceID:            "",
		DisplayCoordinates: nil,
		TrimUser:           nil,
		MediaIds:           nil,
		TweetMode:          "",
	}
	_, _, err := client.Statuses.Update(fmt.Sprintf("Hello @%s , Sorry but I do not understand your message!", tweet.User.ScreenName), statusUpdate)
	if err != nil {
		Logger.Error(fmt.Sprintf("Faild to reply with do not understand, %s", err.Error()))
	}
}

func makeTweetPicAndShare(tweet twitter.Tweet) {
	log.Printf("replyWithScreenShotFor: %s\n", tweet.IDStr)

	filename, err := TweetScreenShot(tweet.User.ScreenName, tweet.IDStr)
	if err != nil {
		log.Fatal(err)
		Logger.Error(fmt.Sprintf("Faild to take a screenshot of the tweet, %s", err.Error()))
	}

	// tweeting with photos is not yet supported in the tweeter sdk library
	// so I'll use only url of the image to be part of the text :/

	statusUpdate := &twitter.StatusUpdateParams{
		Status:             "",
		InReplyToStatusID:  tweet.ID,
		PossiblySensitive:  nil,
		Lat:                nil,
		Long:               nil,
		PlaceID:            "",
		DisplayCoordinates: nil,
		TrimUser:           nil,
		MediaIds:           nil,
		TweetMode:          "",
	}

	filename = makeURL(filename)

	_, _, err2 := client.Statuses.Update(fmt.Sprintf("Hello @%s , Here you are %s", tweet.User.ScreenName, filename), statusUpdate)
	if err2 != nil {
		Logger.Error(fmt.Sprintf("Faild to reply pic tweet, %s", err2.Error()))
	}
}

func makeURL(s string) string {
	return fmt.Sprintf("%s%s", PIC_STORAGE_URL, s)
}
