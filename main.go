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

func init() {
	viper.SetDefault("tag_names", []string{"r3", "realraum"})
	pflag.StringSliceVar(&DebugFlags_, "debug", []string{}, "debug flags e.g. ALL,MADON,MAIN")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.SetEnvPrefix("MBB")
	viper.AutomaticEnv()
}

func main() {
	LogEnable(viper.GetStringSlice("debug")...)
	status_lvl1 := make(chan madon.Status, 20)
	status_lvl2 := make(chan madon.Status, 15)

	tag_names := viper.GetStringSlice("tag_names")

	client, err := madonMustInitClient()
	if err != nil {
		LogMain_.Fatal(err)
	}

	go goSubscribeStreamOfTagNames(client, tag_names, status_lvl1)
	go goFilterStati(client, status_lvl1, status_lvl2, StatusFilterConfig{must_have_visiblity: []string{"public"}, must_have_one_of_tag_names: tag_names, must_be_unmuted: true, must_be_original: true, must_be_followed_by_us: true, must_not_be_sensitive: false})

	// wait on Ctrl-C or sigInt or sigKill
	go func() {
		ctrlc_c := make(chan os.Signal, 1)
		signal.Notify(ctrlc_c, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-ctrlc_c //block until ctrl+c is pressed || we receive SIGINT aka kill -1 || kill
		LogMain_.Println("SIGINT received, exiting gracefully ...")
		os.Exit(0)
	}()

	//goBoostStati(client, status_lvl2)
	goPrintStati(status_lvl2)

	LogMain_.Print("Exiting..")
}
