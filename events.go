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
				filename, hit, media := selectAndDownload(link)

				if filename != "" {
					if media == Video {
						if _, ok := cacheVideo[filename]; ok {
							videos := *cacheVideo[filename]
							album := make(tele.Album, 0, 10)

							for i := 0; i < len(videos); i += 10 {
								// Add photos to album
								for j := 0; j < 10; j++ {
									if i+j < len(videos) {
										album = append(album, videos[i+j])
									}
								}

								err := c.SendAlbum(album, tele.Silent)
								if err != nil {
									lit.Error(err.Error())
								}
							}

							if !hit {
								go saveVideo(cacheVideo[filename])
							}
						}
					} else {
						if _, ok := cacheAlbum[filename]; ok {
							photos := *cacheAlbum[filename]
							album := make(tele.Album, 0, 10)

							for i := 0; i < len(photos); i += 10 {
								// Add photos to album
								for j := 0; j < 10; j++ {
									if i+j < len(photos) {
										album = append(album, photos[i+j])
									}
								}

								_ = c.SendAlbum(album, tele.Silent)
							}

							if !hit {
								go saveAlbum(&photos, filename)
							}

							// Handle audio
							filename, hit = downloadAudio(link)
							if filename != "" {
								err := c.Reply(cacheAudio[filename], tele.Silent)
								if err == nil && !hit {
									go saveAudio(cacheAudio[filename])
								}
							}
						}
					}
				}
			} else {
				// For twitter, we send the same url with only fx appended to it
				if strings.HasPrefix(link, "https://twitter.com") {
					err := c.Reply(strings.Replace(link, "https://twitter.com", "https://vxtwitter.com", 1), tele.Silent)
					if err != nil {
						lit.Error(err.Error())
					}
				} else {
					if strings.HasPrefix(link, "https://x.com") {
						err := c.Reply(strings.Replace(link, "https://x.com", "https://vxtwitter.com", 1), tele.Silent)
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
		text    = c.Query().Text
		results = make([]tele.Result, 0, 1)
	)

	if isValidURL(text) && contains(text, cfg.URLs) {
		text = cleanURL(text)

		filename, hit, media := selectAndDownload(text)

		if filename != "" {
			if media == Video {
				if !hit {
					// Upload videos to channel, so we can send it even in inline mode
					for _, v := range *cacheVideo[filename] {
						_, err := c.Bot().Send(tele.ChatID(cfg.Channel), v)
						if err != nil {
							return err
						}
					}

					go saveVideo(cacheVideo[filename])
				}

				for i, v := range *cacheVideo[filename] {
					results = append(results, &tele.VideoResult{
						Cache: v.FileID,
						Title: "Send video",
						MIME:  "video/mp4",
					})

					results[i].SetResultID(filename)
				}
			} else {
				if _, ok := cacheAlbum[filename]; ok {
					photos := *cacheAlbum[filename]

					for i, p := range photos {
						_, _ = c.Bot().Send(tele.ChatID(cfg.Channel), p)

						results = append(results, &tele.PhotoResult{
							URL:      p.FileURL,
							ThumbURL: p.FileURL,
						})
						results[i].SetResultID(idGen(p.FileURL))
					}

					if !hit {
						go saveAlbum(&photos, filename)
					}
				}
			}
		}
	} else {
		results = append(results, &tele.ArticleResult{
			Title: "Not a valid URL!",
		})

		results[0].SetResultID(text)
	}

	return c.Answer(&tele.QueryResponse{
		Results:   results,
		CacheTime: 86400, // one day
	})
}
