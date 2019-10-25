package main

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

// TweetScreenShot captures a screenshot of a tweet
func TweetScreenShot(username string, tweetID string) ([]byte, error) {

	chromedpContext, cancelCtxt := chromedp.NewContext(context.Background()) // create new tab
	defer cancelCtxt()

	// capture screenShot of an element
	url := fmt.Sprintf("https://twitter.com/%s/status/%s", username, tweetID)
	selector := `document.querySelector("#permalink-overlay-dialog > div.PermalinkOverlay-content > div > div > div.permalink.light-inline-actions.stream-uncapped.original-permalink-page > div.permalink-inner.permalink-tweet-container > div")`
	var buf []byte
	if err := chromedp.Run(chromedpContext, elementScreenshot(url, selector, &buf)); err != nil {
		return buf, err
	}

	return buf, nil
}

// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitVisible(sel, chromedp.ByJSPath),
		chromedp.Sleep(time.Second * 3),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible, chromedp.ByJSPath),
	}
}
