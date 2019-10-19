package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

// MediaUpload URI
const MediaUpload string = "https://upload.twitter.com/1.1/media/upload.json"

// TwitterUpload http client
type TwitterUpload struct {
	client *http.Client
}

// MediaInitResponse properties
type MediaInitResponse struct {
	MediaID          int64  `json:"media_id"`
	MediaIDString    string `json:"media_id_string"`
	ExpiresAfterSecs uint64 `json:"expires_after_secs"`
}

// NewTwitterUpload uploads new screenshots
func NewTwitterUpload(client *http.Client) *TwitterUpload {
	self := &TwitterUpload{}
	self.client = client
	return self
}

// Upload function, it uploads tweet screenshot in a reply on twitter
func (twitterUploader *TwitterUpload) Upload(path string) (int64, error) {
	media, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}

	mediaInitResponse, err := twitterUploader.mediaInit(media)
	if err != nil {
		return 0, err
	}

	mediaID := mediaInitResponse.MediaID

	if twitterUploader.mediaAppend(mediaID, media, path) != nil {
		return 0, err
	}

	if twitterUploader.mediaFinilize(mediaID) != nil {
		return 0, err
	}

	return mediaID, nil
}

func (twitterUploader *TwitterUpload) mediaInit(media []byte) (*MediaInitResponse, error) {
	form := url.Values{}
	form.Add("command", "INIT")
	form.Add("media_type", "image/png")
	form.Add("total_bytes", fmt.Sprint(len(media)))

	req, err := http.NewRequest("POST", MediaUpload, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := twitterUploader.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	var mediaInitResponse MediaInitResponse
	err = json.Unmarshal(body, &mediaInitResponse)

	if err != nil {
		return nil, err
	}

	return &mediaInitResponse, nil
}

func (twitterUploader *TwitterUpload) mediaAppend(mediaID int64, media []byte, path string) error {
	step := 500 * 1024
	for s := 0; s*step < len(media); s++ {
		var body bytes.Buffer
		rangeBegining := s * step
		rangeEnd := (s + 1) * step
		if rangeEnd > len(media) {
			rangeEnd = len(media)
		}

		w := multipart.NewWriter(&body)

		w.WriteField("command", "APPEND")
		w.WriteField("media_id", fmt.Sprint(mediaID))
		w.WriteField("segment_index", fmt.Sprint(s))

		fw, err := w.CreateFormFile("media", path)

		fw.Write(media[rangeBegining:rangeEnd])

		w.Close()

		req, err := http.NewRequest("POST", MediaUpload, &body)

		req.Header.Add("Content-Type", w.FormDataContentType())

		res, err := twitterUploader.client.Do(req)
		if err != nil {
			return err
		}

		resBody, err := ioutil.ReadAll(res.Body)
		_ = resBody
	}

	return nil
}

func (twitterUploader *TwitterUpload) mediaFinilize(mediaID int64) error {
	form := url.Values{}
	form.Add("command", "FINALIZE")
	form.Add("media_id", fmt.Sprint(mediaID))

	req, err := http.NewRequest("POST", MediaUpload, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := twitterUploader.client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	_ = body

	return nil
}
