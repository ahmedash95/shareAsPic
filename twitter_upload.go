package main

import (
	"io/ioutil"
	"fmt"
	"net/http"
	"encoding/json"
	"bytes"
	"mime/multipart"
	"net/url"
	"strings"
)

const MediaUpload string = "https://upload.twitter.com/1.1/media/upload.json"

type TwitterUpload struct {
	client *http.Client
}

type MediaInitResponse struct {
	MediaId int64 `json:"media_id"`
	MediaIdString string `json:"media_id_string"`
	ExpiresAfterSecs uint64 `json:"expires_after_secs"`
}

func NewTwitterUpload(client *http.Client) *TwitterUpload {
	self := &TwitterUpload{}
	self.client = client
	return self
}

func (self *TwitterUpload) Upload(path string) (int64,error) {
	media, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}

	mediaInitResponse, err := self.MediaInit(media)
	if err != nil {
		return 0, err
	}

	mediaId := mediaInitResponse.MediaId

	if self.MediaAppend(mediaId, media,path) != nil {
		return 0, err
	}

	if self.MediaFinilize(mediaId) != nil {
		return 0, err
	}

	return mediaId,nil
}

func (self *TwitterUpload) MediaInit(media []byte) (*MediaInitResponse, error) {
	form := url.Values{}
	form.Add("command", "INIT")
	form.Add("media_type", "image/png")
	form.Add("total_bytes", fmt.Sprint(len(media)))


	req, err := http.NewRequest("POST", MediaUpload, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := self.client.Do(req)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	var mediaInitResponse MediaInitResponse
	err = json.Unmarshal(body, &mediaInitResponse)

	if err != nil {
		return nil, err
	}


	return &mediaInitResponse, nil
}

func (self *TwitterUpload) MediaAppend(mediaId int64, media []byte,path string) error {
	step := 500 * 1024
	for s := 0; s * step < len(media); s++ {
		var body bytes.Buffer
		rangeBegining := s * step
		rangeEnd := (s + 1) * step
		if rangeEnd > len(media) {
			rangeEnd = len(media)
		}


		w := multipart.NewWriter(&body)

		w.WriteField("command", "APPEND")
		w.WriteField("media_id", fmt.Sprint(mediaId))
		w.WriteField("segment_index", fmt.Sprint(s))

		fw, err := w.CreateFormFile("media", path)


		fw.Write(media[rangeBegining:rangeEnd])


		w.Close()

		req, err := http.NewRequest("POST", MediaUpload, &body)

		req.Header.Add("Content-Type", w.FormDataContentType())

		res, err := self.client.Do(req)
		if err != nil {
			return err
		}

		resBody, err := ioutil.ReadAll(res.Body)
		_ = resBody
	}

	return nil
}

func (self *TwitterUpload) MediaFinilize(mediaId int64) error {
	form := url.Values{}
	form.Add("command", "FINALIZE")
	form.Add("media_id", fmt.Sprint(mediaId))

	req, err := http.NewRequest("POST", MediaUpload, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := self.client.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)
	_ = body

	return nil
}
