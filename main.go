package main

import (
	"github.com/bwmarrin/lit"
	"github.com/kkyr/fig"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

const (
	tempFolder = "./temp/"
)

var (
	cfg config
)

func init() {
	lit.LogLevel = lit.LogError

	err := fig.Load(&cfg, fig.File("config.yml"))
	if err != nil {
		lit.Error(err.Error())
		return
	}

	// Set lit.LogLevel to the given value
	switch strings.ToLower(cfg.LogLevel) {
	case "logwarning", "warning":
		lit.LogLevel = lit.LogWarning

	case "loginformational", "informational":
		lit.LogLevel = lit.LogInformational

	case "logdebug", "debug":
		lit.LogLevel = lit.LogDebug
	}

	// Create folders used by the bot
	if _, err = os.Stat(tempFolder); err != nil {
		if err = os.Mkdir(tempFolder, 0755); err != nil {
			lit.Error("Cannot create temp directory, %s", err)
		}
	}
}

func main() {
	// Start HTTP server to serve generated .mp3 files
	http.Handle("/temp/", http.StripPrefix("/temp", http.FileServer(http.Dir(tempFolder))))
	go http.ListenAndServe(":8070", nil)

	// Create bot
	pref := tele.Settings{
		Token:  cfg.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle(tele.OnQuery, inlineQuery)

	b.Handle(tele.OnText, videoDownload)
	b.Handle(tele.OnText, twitterReplacer)

	// Start bot
	lit.Info("videoDownloader is now running")
	b.Start()
}
