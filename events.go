package main

import (
	"github.com/bwmarrin/lit"
	tele "gopkg.in/telebot.v3"
	"strings"
)

func videoDownload(c tele.Context) error {
	for _, e := range c.Message().Entities {
		if e.Type == tele.EntityURL {
			url := cleanURL(c.Message().EntityText(e))

			if contains(url, cfg.URLs) {
				filename, hit := checkAndDownload(url)

				err := c.Reply(cache[filename], tele.Silent)
				if err != nil {
					lit.Error(err.Error())
				}

				if !hit {
					save(cache[filename])
				}
			} else {
				// For twitter, we send the same url with only fx appended to it
				if strings.HasPrefix(url, "https://twitter.com") {
					err := c.Reply(strings.Replace(url, "https://twitter.com", "https://fxtwitter.com", 1), tele.Silent)
					if err != nil {
						lit.Error(err.Error())
					}
				}
			}
		}
	}

	return nil
}

func inlineQuery(c tele.Context) error {
	var (
		results = make([]tele.Result, 1)
		text    = c.Query().Text
	)

	if isValidURL(text) && contains(text, cfg.URLs) {
		text = cleanURL(text)
		filename, hit := checkAndDownload(text)

		// Upload video to channel, so we can send it even in inline mode
		_, _ = c.Bot().Send(tele.ChatID(cfg.Channel), cache[filename])

		if !hit {
			go save(cache[filename])
		}

		// Create result
		results[0] = &tele.VideoResult{
			Cache: cache[filename].FileID,
			Title: "Send video",
			MIME:  "video/mp4",
		}

		results[0].SetResultID(filename)

		// Send video
		return c.Answer(&tele.QueryResponse{
			Results:   results,
			CacheTime: 86400, // one day
		})
	} else {
		results[0] = &tele.ArticleResult{
			Title: "Not a valid instagram URL!",
		}

		results[0].SetResultID(text)

		// Send error
		return c.Answer(&tele.QueryResponse{
			Results:   results,
			CacheTime: 86400, // one day
		})
	}
}
