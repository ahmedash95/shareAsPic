package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

const ProcessedTweets = "processed_tweets"

var RepliesSet = []string{
	"مساء الخير يا %s اتفضل يا زعيم",
	"انت تؤمر يا %s",
	"طلباتك اوامر يا  %s",
	"لو عجبتك يا %s اعمل فولو بقي",
	"اتفضل يا %s و ادعيلنا دعوتين حلوين",
	"ازيك يا %s اتفضل خلصتهالك",
}

var client *twitter.Client
var httpClient *http.Client
var TwitterUploadClient *TwitterUpload

func initTwitterClient() {
	config := oauth1.NewConfig(TWITTER_API_KEY, TWITTER_API_SECRET)
	token := oauth1.NewToken(TWITTER_ACCESS_TOKEN_KEY, TWITTER_ACCESS_TOKEN_SECRET)
	httpClient = config.Client(oauth1.NoContext, token)
	client = twitter.NewClient(httpClient)
	TwitterUploadClient = NewTwitterUpload(httpClient)
}

func processTweet(tweet twitter.Tweet) {
	if tweetProcessedBefore(tweet) {
		logAndPrint(fmt.Sprintf("Tweet proccessed before %s/status/%s", tweet.InReplyToScreenName, tweet.InReplyToStatusIDStr))
		return
	}
	// let's make sure it has "share this" in the string
	if !strings.Contains(strings.ToLower(tweet.Text), "@shareaspic") || !strings.Contains(strings.ToLower(tweet.Text), "share this") {
		return
	}

	makeTweetPicAndShare(tweet)
}

func tweetProcessedBefore(tweet twitter.Tweet) bool {
	result, _ := redisClient.SAdd(ProcessedTweets, tweet.ID).Result()
	var processed bool
	processed = result == 0
	if !processed {
		log.Printf("new tweet to process: %s\n", tweet.IDStr)
	}

	return processed
}

func makeTweetPicAndShare(tweet twitter.Tweet) {
	logAndPrint(fmt.Sprintf("prepare replyWithScreenShotFor: %s\n", tweet.IDStr))

	logAndPrint("taking a screenshot")
	filename, err := TweetScreenShot(tweet.InReplyToScreenName, tweet.InReplyToStatusIDStr)
	if err != nil {
		logAndPrint(fmt.Sprintf("Faild to take a screenshot of the tweet, %s", err.Error()))
		return
	}
	logAndPrint("screenshot has been taken successfully")

	logAndPrint(fmt.Sprintf("replying to %s (%s) for reply to %s/status/%s", tweet.User.ScreenName, tweet.IDStr, tweet.InReplyToScreenName, tweet.InReplyToStatusIDStr))

	filePath := fmt.Sprintf("%s%s", PIC_STORAGE_PATH, filename)

	logAndPrint("upload photo")
	mediaId, err := TwitterUploadClient.Upload(filePath)
	logAndPrint(fmt.Sprintf("photo has been uploaded: %d", mediaId))

	statusUpdate := &twitter.StatusUpdateParams{
		Status:             "",
		InReplyToStatusID:  tweet.ID,
		PossiblySensitive:  nil,
		Lat:                nil,
		Long:               nil,
		PlaceID:            "",
		DisplayCoordinates: nil,
		TrimUser:           nil,
		MediaIds:           []int64{mediaId},
		TweetMode:          "",
	}

	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(RepliesSet)

	_, _, err2 := client.Statuses.Update(fmt.Sprintf(RepliesSet[n], tweet.User.ScreenName), statusUpdate)
	if err2 != nil {
		logAndPrint(fmt.Sprintf("Faild to reply pic tweet, %s", err2.Error()))
	}

	logAndPrint(fmt.Sprintf("replied With screenshot for: %s\n", tweet.IDStr))
}
