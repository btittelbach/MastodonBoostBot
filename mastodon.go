package main

import (
	"fmt"
	"io"
	"os"

	"github.com/McKael/madon"
	"github.com/spf13/viper"
)

/// code adapted from madonctl by McKael -- https://github.com/McKael/madonctl
/// Kudos!!
func madonMustInitClient() (client *madon.Client, err error) {

	appName := viper.GetString("AppName")
	instanceURL := viper.GetString("instance")
	appID := viper.GetString("app_id")
	appSecret := viper.GetString("app_secret")

	if instanceURL == "" {
		LogMadon_.Fatalln("madonInitClient:", "no instance provided")
	}

	LogMadon_.Println("madonInitClient:Instance: '%s'", instanceURL)

	if appID != "" && appSecret != "" {
		// We already have an app key/secret pair
		client, err = madon.RestoreApp(appName, instanceURL, appID, appSecret, nil)
		if err != nil {
			return
		}
		// Check instance
		if _, err = client.GetCurrentInstance(); err != nil {
			LogMadon_.Fatalln("madonInitClient:", err, "could not connect to server with provided app ID/secret")
			return
		}
		LogMadon_.Println("madonInitClient:", "Using provided app ID/secret")
		return
	}

	if appID != "" || appSecret != "" {
		LogMadon_.Fatalln("madonInitClient:", "Warning: provided app id/secrets incomplete -- registering again")
	}

	LogMadon_.Println("madonInitClient:", "Registered new application.")
	return
}

/// code adapted from madonctl by McKael -- https://github.com/McKael/madonctl
/// Kudos!!
func goSubscribeStreamOfTagNames(client *madon.Client, hashTagList []string, statusOutChan chan<- madon.Status) {
	streamName := "hashtag"
	evChan := make(chan madon.StreamEvent, 10)
	var param string
	stop := make(chan bool)
	done := make(chan bool)
	var err error

	if len(hashTagList) <= 1 { // Usual case: Only 1 stream
		err = client.StreamListener(streamName, param, evChan, stop, done)
	} else { // Several streams
		n := len(hashTagList)
		tagEvCh := make([]chan madon.StreamEvent, n)
		tagDoneCh := make([]chan bool, n)
		for i, t := range hashTagList {
			LogMadon_.Println("goSubscribeStreamOfTagNames: Launching listener for tag '%s'", t)
			tagEvCh[i] = make(chan madon.StreamEvent)
			tagDoneCh[i] = make(chan bool)
			e := client.StreamListener(streamName, t, tagEvCh[i], stop, tagDoneCh[i])
			if e != nil {
				if i > 0 { // Close previous connections
					close(stop)
				}
				err = e
				break
			}
			// Forward events to main ev channel
			go func(i int) {
				for {
					select {
					case _, ok := <-tagDoneCh[i]:
						if !ok { // end of streaming for this tag
							done <- true
							return
						}
					case ev := <-tagEvCh[i]:
						evChan <- ev
					}
				}
			}(i)
		}
	}

	if err != nil {
		LogMadon_.Fatalln("goSubscribeStreamOfTagNames:", err.Error())
	}

LISTENSTREAM:
	for {
		select {
		case v, ok := <-done:
			if !ok || v == true { // done is closed, end of streaming
				break LISTENSTREAM
			}
		case ev := <-evChan:
			switch ev.Event {
			case "error":
				if ev.Error != nil {
					if ev.Error == io.ErrUnexpectedEOF {
						LogMadon_.Println("goSubscribeStreamOfTagNames:", "The stream connection was unexpectedly closed")
						continue
					}
					LogMadon_.Println("goSubscribeStreamOfTagNames:", "Error event: [%s] %s", ev.Event, ev.Error)
					continue
				}
				LogMadon_.Println("goSubscribeStreamOfTagNames:", "Event: [%s]", ev.Event)
			case "update":
				s := ev.Data.(madon.Status)
				statusOutChan <- s
			case "notification", "delete":
				continue
			default:
				LogMadon_.Println("goSubscribeStreamOfTagNames:", "Unhandled event: [%s] %T", ev.Event, ev.Data)
			}
		}
	}
	close(stop)
	close(evChan)
	close(statusOutChan)
	if err != nil {
		LogMadon_.Println("goSubscribeStreamOfTagNames: Error: %s", err.Error())
		os.Exit(1)
	}
}

func getRelation(client *madon.Client, accID int64) (madon.Relationship, error) {
	relationshiplist, err := client.GetAccountRelationships([]int64{accID})
	if err != nil {
		return madon.Relationship{}, err
	}
	if len(relationshiplist) == 0 {
		return madon.Relationship{}, fmt.Errorf("Unknown error, empty list")
	}
	return relationshiplist[0], nil
}

func goBoostStati(client *madon.Client, stati_chan <-chan madon.Status) {
	for status := range stati_chan {
		LogMadon_.Printf("Boosting Status with ID %d published by %s\n", status.ID, status.Account.Username)
		client.ReblogStatus(status.ID)
	}
}

func goPrintStati(stati_chan <-chan madon.Status) {
	for status := range stati_chan {
		fmt.Printf("%+v\n", status)
	}
}
