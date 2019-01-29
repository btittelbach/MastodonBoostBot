// (c) Bernhard Tittelbach, 2019

package main

import "os"
import "log"

type NullWriter struct{}

func (n *NullWriter) Write(p []byte) (int, error) { return len(p), nil }

var (
	LogMadon_ *log.Logger
	LogMain_  *log.Logger
)

func init() {
	if LogMain_ == nil {
		LogMadon_ = log.New(&NullWriter{}, "", 0)
		LogMain_ = log.New(&NullWriter{}, "", 0)
	}
}

func LogEnable(logtypes ...string) {
	LogMadon_ = log.New(&NullWriter{}, "", 0)
	LogMain_ = log.New(&NullWriter{}, "", 0)
	for _, logtype := range logtypes {
		switch logtype {
		case "MADON":
			LogMadon_ = log.New(os.Stderr, logtype+" ", log.LstdFlags)
		case "MAIN":
			LogMain_ = log.New(os.Stderr, logtype+" ", log.LstdFlags)
		case "ALL":
			LogMadon_ = log.New(os.Stderr, logtype+" ", log.LstdFlags)
			LogMain_ = log.New(os.Stderr, "MAIN"+" ", log.LstdFlags)
		}
	}
}
