package main

import (
	"bytes"
	"io"
	"net/url"
	"strconv"

	"encoding/base64"

	"github.com/ChimeraCoder/anaconda"
	"github.com/spf13/viper"
)

/////////////
/// Twitter
/////////////

func initTwitterClient() *anaconda.TwitterApi {
	if !viper.IsSet("cctweet_access_token") || !viper.IsSet("cctweet_consumer_key") {
		return nil
	}
	return anaconda.NewTwitterApiWithCredentials(
		viper.GetString("cctweet_access_token"),
		viper.GetString("cctweet_access_secret"),
		viper.GetString("cctweet_consumer_key"),
		viper.GetString("cctweet_consumer_secret"))
}

func sendTweet(client *anaconda.TwitterApi, post string, media_ids []string) error {
	var err error
	v := url.Values{}
	v.Set("status", post)
	for _, mid := range media_ids {
		v.Add("media_ids", mid)
	}
	// log.Println("sendTweet", post, v)
	_, err = client.PostTweet(post, v)
	return err
}

func getImageForTweet(client *anaconda.TwitterApi, imagebuffer io.Reader) (string, error) {
	b64buf := new(bytes.Buffer)
	b64encoder := base64.NewEncoder(base64.StdEncoding, b64buf)
	io.Copy(b64encoder, imagebuffer)
	tmedia, err := client.UploadMedia(b64buf.String())
	if err == nil {
		return strconv.FormatInt(tmedia.MediaID, 10), err
	} else {
		return "", err
	}
}
