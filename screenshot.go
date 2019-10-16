package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/chromedp/chromedp"
)

var chromedpContext context.Context

func initChromedpClient() {
	var cancelCtxt context.CancelFunc

	chromedpContext, cancelCtxt = chromedp.NewContext(context.Background()) // create new tab
	_ = cancelCtxt
}

func TweetScreenShot(username string, tweetId string) (string, error) {
	// capture screenShot of an element
	fname := fmt.Sprintf("%s-%s.png", username, tweetId)
	url := fmt.Sprintf("https://twitter.com/%s/status/%s", username, tweetId)

	var buf []byte
	if err := chromedp.Run(chromedpContext, elementScreenshot(url, `document.querySelector("#permalink-overlay-dialog > div.PermalinkOverlay-content > div > div > div.permalink.light-inline-actions.stream-uncapped.original-permalink-page > div.permalink-inner.permalink-tweet-container > div")`, &buf)); err != nil {
		return "", err
	}
	fmt.Printf("write pic to path %s\n", fmt.Sprintf("%s/%s", PIC_STORAGE_PATH, fname))
	if err := ioutil.WriteFile(fmt.Sprintf("%s/%s", PIC_STORAGE_PATH, fname), buf, 0755); err != nil {
		return "", err
	}
	return fname, nil
}

// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitVisible(sel, chromedp.ByJSPath),
		//chromedp.Evaluate("document.getElementsByClassName('follow-bar')[0].remove();", &res),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible, chromedp.ByJSPath),
	}
}
