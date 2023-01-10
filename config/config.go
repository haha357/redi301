package config

import (
	"flag"
	"fmt"
	"github.com/gookit/color"
	"github.com/sirupsen/logrus"
	"net/url"
	"os"
	"strings"
)

var (
	HttpAddr string
	Target   string
	HttpPort string
)

const (
	LogLevel = logrus.DebugLevel
	Header   = `=========================================================
Welcome to use this application.
If you like this application, you can donate to the author.
USDT(Trc20) wallet: TB8meT4Pm9KFXRJ8SNCfxx4yBGPbk3Ekip
Thank you so much.
=========================================================`
)

func Init() {
	flag.Usage = func() {
		color.Redln(Header)
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&HttpAddr, "a", "0.0.0.0:80", "The listen address and port.")
	flag.StringVar(&Target, "t", "https://www.microsoft.com", "The prefix of target redirect url.")
	flag.Parse()
	ipAndPort := strings.Split(HttpAddr, ":")
	if len(ipAndPort) != 2 {
		logrus.Fatalf("http listen address error...")
	} else {
		HttpPort = ipAndPort[1]
	}
	if _, err := url.Parse(Target); err != nil {
		logrus.Fatalf("prefix of target redirect url error...")
	}
}
