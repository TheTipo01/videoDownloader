package main

import (
	"github.com/bwmarrin/lit"
	tele "gopkg.in/telebot.v3"
	"strings"
)

func videoDownload(c tele.Context) error {
	for _, e := range c.Message().Entities {
		if e.Type == tele.EntityURL {
			if contains(e.URL, cfg.URLs) {

				filename := downloadVideo(e.URL)
				err := c.Reply(&tele.Video{File: tele.FromURL(cfg.Host + filename), FileName: filename, MIME: "video/mp4"}, tele.Silent)
				if err != nil {
					lit.Error(err.Error())
				}
			} else {
				// For twitter we send the same url with only fx appended to it
				if strings.HasPrefix(e.URL, "https://twitter.com") {
					err := c.Reply(strings.Replace(e.URL, "https://twitter.com", "https://fxtwitter.com", 1), tele.Silent)
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
		fileName := downloadVideo(text)

		// Create result
		results[0] = &tele.VideoResult{
			URL:      cfg.Host + fileName,
			Title:    "Send video",
			MIME:     "video/mp4",
			ThumbURL: cfg.Host + "icon.jpg",
		}

		results[0].SetResultID(fileName)

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
