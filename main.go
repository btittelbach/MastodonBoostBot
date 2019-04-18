/// (c) Bernhard Tittelbach, 2019 - MIT License

package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/McKael/madon"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var DebugFlags_ []string
var RegisterAppForTwitterUser_ bool

func init() {
	viper.SetDefault("tag_names", []string{"r3", "realraum"})
	pflag.StringSliceVar(&DebugFlags_, "debug", []string{}, "debug flags e.g. ALL,MADON,MAIN")
	pflag.BoolVar(&RegisterAppForTwitterUser_, "starttwitteroauth", false, "oauth register this app with your twitter user")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.SetEnvPrefix("MBB")
	viper.AutomaticEnv()
	viper.SetDefault("filterout_sensitive", false)
	viper.SetDefault("filterout_reboosts", true)
	viper.SetDefault("filterfor_accounts_we_follow", true)
}

func goSplitChannel(in <-chan madon.Status, out1, out2 chan<- madon.Status) {
	for status := range in {
		out1 <- status
		out2 <- status
	}
}

func main() {
	LogEnable(viper.GetStringSlice("debug")...)

	if RegisterAppForTwitterUser_ {
		oauthAppWithTwitterForUser(viper.GetString("cctweet_consumer_key"), viper.GetString("cctweet_consumer_secret"))
		os.Exit(0)
	}

	status_lvl1 := make(chan madon.Status, 20)
	status_lvl2 := make(chan madon.Status, 15)

	tag_names := viper.GetStringSlice("tag_names")

	client, err := madonMustInitClient()
	if err != nil {
		LogMain_.Fatal(err)
	}

	birdclient := initTwitterClient()

	go goSubscribeStreamOfTagNames(client, tag_names, status_lvl1)
	go goFilterStati(client, status_lvl1, status_lvl2, StatusFilterConfig{must_have_visiblity: []string{"public"}, must_have_one_of_tag_names: tag_names, must_be_unmuted: true, must_be_original: viper.GetBool("filterout_reboosts"), must_be_followed_by_us: viper.GetBool("filterfor_accounts_we_follow"), must_not_be_sensitive: viper.GetBool("filterout_sensitive")})

	if birdclient != nil {
		status_lvl3_twitter := make(chan madon.Status, 15)
		status_lvl3_boost := make(chan madon.Status, 15)
		go goSplitChannel(status_lvl2, status_lvl3_boost, status_lvl3_twitter)
		go goTweetStati(client, birdclient, status_lvl3_twitter)
		go goBoostStati(client, status_lvl3_boost)
	} else {
		go goBoostStati(client, status_lvl2)
		// goPrintStati(status_lvl2)
	}

	// wait on Ctrl-C or sigInt or sigKill
	{
		ctrlc_c := make(chan os.Signal, 1)
		signal.Notify(ctrlc_c, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-ctrlc_c //block until ctrl+c is pressed || we receive SIGINT aka kill -1 || kill
		LogMain_.Println("SIGINT received, exiting gracefully ...")
	}
}
