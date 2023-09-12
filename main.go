package main

import (
	"database/sql"
	"github.com/bwmarrin/lit"
	"github.com/kkyr/fig"
	"log"
	_ "modernc.org/sqlite"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

var (
	cfg   config
	cache map[string]*tele.Video
	db    *sql.DB
)

func init() {
	lit.LogLevel = lit.LogError

	err := fig.Load(&cfg, fig.File("config.yml"), fig.Dirs(".", "./data"))
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

	cache = make(map[string]*tele.Video)

	// Open database connection
	db, err = sql.Open("sqlite", "./data/cache.db")
	if err != nil {
		lit.Error("Error opening db connection, %s", err)
		return
	}

	execQuery(cacheTable)
	load()
}

func main() {
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

	// Start bot
	lit.Info("videoDownloader is now running")
	b.Start()

	_ = db.Close()
}
