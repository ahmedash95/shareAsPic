package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func TweetScreenShot(username string, tweetId string) (string, error) {

	chromedpContext, cancelCtxt := chromedp.NewContext(context.Background()) // create new tab
	defer cancelCtxt()

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

func TweetSendReply(userScreenName, tweetID, message string) error {
	chromedpContext, cancelCtxt := chromedp.NewContext(context.Background()) // create new tab
	defer cancelCtxt()

	var buf []byte
	if err := chromedp.Run(chromedpContext, replyToTweet(userScreenName, tweetID, message, &buf)); err != nil {
		return err
	}

	if err2 := ioutil.WriteFile("./pics/fullScreenshot-2222.png", buf, 0644); err2 != nil {
		log.Fatal(err2)
	}

	return nil
}

// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitVisible(sel, chromedp.ByJSPath),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible, chromedp.ByJSPath),
	}
}

func replyToTweet(userScreen, tweetId, message string, res *[]byte) chromedp.Tasks {
	replyIcon := "document.querySelector('[aria-label=\"Reply\"]')"
	replyTextAreaPlaceHolder := "document.querySelector(\"#react-root > div > div > div.r-1d2f490.r-u8s1d.r-zchlnj.r-ipm5af.r-184en5c > div > div > div > div > div.css-1dbjc4n.r-1habvwh.r-18u37iz.r-1pi2tsx.r-1777fci.r-1xcajam.r-ipm5af.r-g6jmlv > div.css-1dbjc4n.r-t23y2h.r-1wbh5a2.r-rsyp9y.r-1pjcn9w.r-htvplk.r-1udh08x.r-1potc6q > div > div.css-1dbjc4n.r-16y2uox.r-1wbh5a2.r-1jgb5lz.r-1ye8kvj.r-13qz1uu > div > div > div:nth-child(2) > div > div > div > div.css-1dbjc4n.r-1iusvr4.r-46vdb2.r-15d164r.r-9cviqr.r-bcqeeo.r-1bylmt5.r-13tjlyg.r-7qyjyx.r-1ftll1t > div.css-1dbjc4n.r-184en5c > div > div > div > div > div > div > div > div > div.css-901oao.r-hkyrab.r-6koalj.r-16y2uox.r-1qd0xha.r-1b6yd1w.r-16dba41.r-ad9z0x.r-bcqeeo.r-qvutc0 > div > div > div > div.DraftEditor-editorContainer > div > div > div > div\")"
	replyTextAreaInput := "document.querySelector(\"#react-root > div > div > div.r-1d2f490.r-u8s1d.r-zchlnj.r-ipm5af.r-184en5c > div:nth-child(1) > div > div > div > div.css-1dbjc4n.r-1habvwh.r-18u37iz.r-1pi2tsx.r-1777fci.r-1xcajam.r-ipm5af.r-g6jmlv > div.css-1dbjc4n.r-t23y2h.r-1wbh5a2.r-rsyp9y.r-1pjcn9w.r-htvplk.r-1udh08x.r-1potc6q > div > div.css-1dbjc4n.r-16y2uox.r-1wbh5a2.r-1jgb5lz.r-1ye8kvj.r-13qz1uu > div > div > div:nth-child(2) > div > div > div > div.css-1dbjc4n.r-1iusvr4.r-46vdb2.r-15d164r.r-9cviqr.r-bcqeeo.r-1bylmt5.r-13tjlyg.r-7qyjyx.r-1ftll1t > div.css-1dbjc4n.r-184en5c > div > div > div > div > div > div > div > div > div.css-901oao.r-hkyrab.r-6koalj.r-16y2uox.r-1qd0xha.r-1b6yd1w.r-16dba41.r-ad9z0x.r-bcqeeo.r-qvutc0 > div > div > div > div.DraftEditor-editorContainer > div\")"
	replyButton := "document.querySelector(\"#react-root > div > div > div.r-1d2f490.r-u8s1d.r-zchlnj.r-ipm5af.r-184en5c > div:nth-child(1) > div > div > div > div.css-1dbjc4n.r-1habvwh.r-18u37iz.r-1pi2tsx.r-1777fci.r-1xcajam.r-ipm5af.r-g6jmlv > div.css-1dbjc4n.r-t23y2h.r-1wbh5a2.r-rsyp9y.r-1pjcn9w.r-htvplk.r-1udh08x.r-1potc6q > div > div.css-1dbjc4n.r-16y2uox.r-1wbh5a2.r-1jgb5lz.r-1ye8kvj.r-13qz1uu > div > div > div:nth-child(2) > div > div > div > div.css-1dbjc4n.r-1iusvr4.r-46vdb2.r-15d164r.r-9cviqr.r-bcqeeo.r-1bylmt5.r-13tjlyg.r-7qyjyx.r-1ftll1t > div:nth-child(2) > div > div > div:nth-child(2) > div.css-18t94o4.css-1dbjc4n.r-urgr8i.r-42olwf.r-sdzlij.r-1phboty.r-rs99b7.r-1w2pmg.r-1n0xq6e.r-1vuscfd.r-1dhvaqw.r-1fneopy.r-o7ynqc.r-6416eg.r-lrvibr\")"

	_ = replyIcon
	_ = replyTextAreaPlaceHolder
	_ = replyTextAreaInput
	_ = replyButton

	replyTweetActions := []chromedp.Action{
		chromedp.Navigate(fmt.Sprintf("https://mobile.twitter.com/%s/status/%s", userScreen, tweetId)),
		chromedp.Sleep(time.Second * 3),
		chromedp.WaitVisible(replyIcon, chromedp.ByJSPath),
		chromedp.Click(replyIcon, chromedp.ByJSPath),
		chromedp.WaitVisible(replyTextAreaPlaceHolder, chromedp.ByJSPath),
		chromedp.SendKeys(replyTextAreaInput, message, chromedp.ByJSPath),
		chromedp.Click(replyButton, chromedp.ByJSPath),
		chromedp.Sleep(time.Second * 3),
		fullScreenShot(res),
	}
	_ = replyTweetActions
	actions := []chromedp.Action{}
	actions = append(actions, loginFirst(res)...)
	actions = append(actions, replyTweetActions...)

	return actions
}

func loginFirst(res *[]byte) []chromedp.Action {

	loginButtonJSPath := "document.querySelector(\"#react-root > div > div > div.css-1dbjc4n.r-1pi2tsx.r-13qz1uu.r-417010 > main > div > div > form > div > div:nth-child(8) > div\")"
	emailInputJSPath := "document.querySelector(\"#react-root > div > div > div.css-1dbjc4n.r-1pi2tsx.r-13qz1uu.r-417010 > main > div > div > form > div > div:nth-child(6) > label > div.css-1dbjc4n.r-18u37iz.r-16y2uox.r-1wbh5a2.r-1udh08x > div > input\")"
	passwordInputJSPath := "document.querySelector(\"#react-root > div > div > div.css-1dbjc4n.r-1pi2tsx.r-13qz1uu.r-417010 > main > div > div > form > div > div:nth-child(7) > label > div.css-1dbjc4n.r-18u37iz.r-16y2uox.r-1wbh5a2.r-1udh08x > div > input\")"

	_ = loginButtonJSPath
	_ = emailInputJSPath
	_ = passwordInputJSPath
	return []chromedp.Action{
		chromedp.Navigate("https://mobile.twitter.com/login"),
		chromedp.WaitVisible(loginButtonJSPath, chromedp.ByJSPath),
		chromedp.SendKeys(emailInputJSPath, TWITTER_USER, chromedp.ByJSPath),
		chromedp.SendKeys(passwordInputJSPath, TWITTER_PASS, chromedp.ByJSPath),
		chromedp.Click(loginButtonJSPath, chromedp.ByJSPath),
		chromedp.Sleep(time.Second * 3),
	}
}

func fullScreenShot(res *[]byte) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		// get layout metrics
		_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
		if err != nil {
			return err
		}

		width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

		// force viewport emulation
		err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
			WithScreenOrientation(&emulation.ScreenOrientation{
				Type:  emulation.OrientationTypePortraitPrimary,
				Angle: 0,
			}).
			Do(ctx)
		if err != nil {
			return err
		}

		// capture screenshot
		*res, err = page.CaptureScreenshot().
			WithQuality(90).
			WithClip(&page.Viewport{
				X:      contentSize.X,
				Y:      contentSize.Y,
				Width:  contentSize.Width,
				Height: contentSize.Height,
				Scale:  1,
			}).Do(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}
