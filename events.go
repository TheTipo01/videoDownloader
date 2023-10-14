package main

import (
	"github.com/bwmarrin/lit"
	tele "gopkg.in/telebot.v3"
	"strings"
)

func videoDownload(c tele.Context) error {
	for _, e := range c.Message().Entities {
		if e.Type == tele.EntityURL {
			link := cleanURL(c.Message().EntityText(e))

			if contains(link, cfg.URLs) {
				// Use the downloader to get videos and albums from tiktok
				if strings.Contains(link, "tiktok.com") {
					filename, hit, media := downloadTikTok(link)
					if filename != "" {
						if media == Video {
							err := c.Reply(cacheVideo[filename], tele.Silent)
							if err == nil {
								if !hit {
									go saveVideo(cacheVideo[filename])
									continue
								}
							}
						} else {
							if _, ok := cacheAlbum[filename]; ok {
								var err error

								photos := *cacheAlbum[filename]
								album := make(tele.Album, 0, 10)

								for i := 0; i < len(photos); i += 10 {
									// Add photos to album
									for j := 0; j < 10; j++ {
										if i+j < len(photos) {
											album = append(album, photos[i+j])
										}
									}

									err = c.SendAlbum(album, tele.Silent)
									if err != nil {
										break
									}
								}

								if err == nil {
									if !hit {
										go saveAlbum(&photos, filename)
										continue
									}
								}
							}
						}
					}
				}

				filename, hit := downloadYtDlp(link)

				err := c.Reply(cacheVideo[filename], tele.Silent)
				if err == nil && !hit {
					go saveVideo(cacheVideo[filename])
				}
			} else {
				// For twitter, we send the same url with only fx appended to it
				if strings.HasPrefix(link, "https://twitter.com") {
					err := c.Reply(strings.Replace(link, "https://twitter.com", "https://fxtwitter.com", 1), tele.Silent)
					if err != nil {
						lit.Error(err.Error())
					}
				} else {
					if strings.HasPrefix(link, "https://x.com") {
						err := c.Reply(strings.Replace(link, "https://x.com", "https://fixupx.com", 1), tele.Silent)
						if err != nil {
							lit.Error(err.Error())
						}
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
		filename, hit := downloadYtDlp(text)

		// Upload video to channel, so we can send it even in inline mode
		_, err := c.Bot().Send(tele.ChatID(cfg.Channel), cacheVideo[filename])
		if err != nil {
			return err
		}

		if !hit {
			go saveVideo(cacheVideo[filename])
		}

		// Create result
		results[0] = &tele.VideoResult{
			Cache: cacheVideo[filename].FileID,
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
