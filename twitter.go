package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"

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
	LogMadon_.Println("sendTweet", post, v)
	_, err = client.PostTweet(post, v)
	return err
}

func getImageForTweet(client *anaconda.TwitterApi, imagebuffer io.Reader) (string, error) {
	b64buf := new(bytes.Buffer)
	b64encoder := base64.NewEncoder(base64.StdEncoding, b64buf)
	io.Copy(b64encoder, imagebuffer)
	tmedia, err := client.UploadMedia(b64buf.String())
	if err == nil {
		LogMadon_.Println("getImageForTweet OK:", tmedia.MediaID)
		return strconv.FormatInt(tmedia.MediaID, 10), err
	} else {
		LogMadon_.Println("getImageForTweet ERROR:", err)
		return "", err
	}
}

func oauthAppWithTwitterForUser(consumerkey, consumersecret string) {
	if !viper.IsSet("cctweet_consumer_secret") || !viper.IsSet("cctweet_consumer_key") {
		panic("consumer keys not set")
	}
	anaconda.SetConsumerKey(consumerkey)
	anaconda.SetConsumerSecret(consumersecret)
	tapi := anaconda.NewTwitterApi("", "")
	url_for_user, temp_credentials, err := tapi.AuthorizationURL("oob")
	if err != nil {
		panic(err)
	}
	fmt.Println("Please authorize this App here: ", url_for_user)
	fmt.Println("Afterwards you will be shown a PIN number.")
	fmt.Println("Enter that PIN here and press Enter to continue: ")
	reader := bufio.NewReader(os.Stdin)
	userpin, _ := reader.ReadString('\n')
	userpin = strings.TrimSpace(userpin)
	tw_credentials, _, err := tapi.GetCredentials(temp_credentials, userpin)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\nPut these into your environment / configuration:\n")
	fmt.Printf("MBB_CCTWEET_CONSUMER_KEY=%s\n", consumerkey)
	fmt.Printf("MBB_CCTWEET_CONSUMER_SECRET=%s\n", consumersecret)
	fmt.Printf("MBB_CCTWEET_ACCESS_TOKEN=%s\n", tw_credentials.Token)
	fmt.Printf("MBB_CCTWEET_ACCESS_SECRET=%s\n", tw_credentials.Secret)
}
