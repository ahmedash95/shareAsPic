package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/ahmedash95/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

const processedTweets = "processed_tweets"

var repliesSet = []string{
	"مساء الخير يا @%s اتفضل يا زعيم",
	"انت تؤمر يا @%s",
	"طلباتك اوامر يا  @%s",
	"لو عجبتك يا @%s اعمل فولو بقي",
	"اتفضل يا @%s و ادعيلنا دعوتين حلوين",
	"ازيك يا @%s اتفضل خلصتهالك",
	"مرحبا @%s شاركها بسهوله الان",
	"مرحبتين @%s اليك ما طلبت",
	"هلا @%s بس تؤمر و احنا ننفذ",
	"هلا @%s في خدمتك",
	"Hey @%s , here you go",
	"Hello @%s , screenshot is ready",
	"Hi @%s , Now you can easily share",
	"Dear @%s , Share now",
}

var client *twitter.Client
var httpClient *http.Client

func initTwitterClient() {
	config := oauth1.NewConfig(twitterAPIKey, twitterAPISECRET)
	token := oauth1.NewToken(twitterAccessTokenKey, twitterAccessTokenSecret)
	httpClient = config.Client(oauth1.NoContext, token)
	client = twitter.NewClient(httpClient)
}

func processTweet(tweet twitter.Tweet) {
	if tweetProcessedBefore(tweet) {
		logAndPrint(fmt.Sprintf("Tweet proccessed before %s/status/%s", tweet.InReplyToScreenName, tweet.InReplyToStatusIDStr))
		return
	}
	if !validMessage(tweet.Text) {
		return
	}
	makeTweetPicAndShare(tweet)
}

func validMessage(tweetText string) bool {
	keywords := []string{"share", "screen shot", "screenshot", "shot", "picture", "tweet"}

	if !strings.Contains(strings.ToLower(tweetText), "@shareaspic") {
		return false
	}

	for i := 0; i < len(keywords); i++ {
		if strings.Contains(strings.ToLower(tweetText), keywords[i]) {
			return true
		}
	}
	return false
}

func tweetProcessedBefore(tweet twitter.Tweet) bool {
	result, _ := redisClient.SAdd(processedTweets, tweet.ID).Result()
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
	screenshot, err := TweetScreenShot(tweet.InReplyToScreenName, tweet.InReplyToStatusIDStr)
	if err != nil {
		logAndPrint(fmt.Sprintf("Faild to take a screenshot of the tweet, %s", err.Error()))
		return
	}
	logAndPrint("screenshot has been taken successfully")

	logAndPrint(fmt.Sprintf("replying to %s (%s) for reply to %s/status/%s", tweet.User.ScreenName, tweet.IDStr, tweet.InReplyToScreenName, tweet.InReplyToStatusIDStr))

	logAndPrint("upload photo")
	mediaUpload := &twitter.MediaUploadParams{
		File:     screenshot,
		MimeType: "image/png",
	}

	media, _, err := client.Media.Upload(mediaUpload)
	if err != nil {
		log.Fatal(err)
	}
	logAndPrint(fmt.Sprintf("photo has been uploaded: %d", media.MediaID))

	statusUpdate := &twitter.StatusUpdateParams{
		Status:             "",
		InReplyToStatusID:  tweet.ID,
		PossiblySensitive:  nil,
		Lat:                nil,
		Long:               nil,
		PlaceID:            "",
		DisplayCoordinates: nil,
		TrimUser:           nil,
		MediaIds:           []int64{media.MediaID},
		TweetMode:          "",
	}

	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(repliesSet)

	_, _, err2 := client.Statuses.Update(fmt.Sprintf(repliesSet[n], tweet.User.ScreenName), statusUpdate)
	if err2 != nil {
		logAndPrint(fmt.Sprintf("Faild to reply pic tweet, %s", err2.Error()))
	}

	logAndPrint(fmt.Sprintf("replied With screenshot for: %s\n", tweet.IDStr))
}
