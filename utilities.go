package main

import (
	"crypto/sha256"
	"encoding/base32"
	"encoding/json"
	"github.com/bwmarrin/lit"
	tele "gopkg.in/telebot.v3"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
)

type Media uint8

const (
	Video Media = iota
	Album
)

// Checks if a string is a valid URL
func isValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	return err == nil
}

// idGen returns the first 11 characters of the SHA1 hash for the given link
func idGen(link string) string {
	h := sha256.New()
	h.Write([]byte(link))

	return strings.ToLower(base32.HexEncoding.EncodeToString(h.Sum(nil))[0:11])
}

func contains(str string, checkFor []string) bool {
	for _, s := range checkFor {
		if strings.Contains(str, s) {
			return true
		}
	}

	return false
}

func downloadYtDlp(link string) (string, bool) {
	hit := true

	filename := idGen(link) + ".mp4"

	if _, ok := cacheVideo[filename]; !ok {
		hit = false

		// Gets info about songs
		info := exec.Command("yt-dlp", "--ignore-errors", "-q", "--no-warnings", "-j", "-a", "-")
		info.Stdin = strings.NewReader(link)

		out, err := info.CombinedOutput()
		if err != nil {
			lit.Error(err.Error())
			return "", hit
		}

		splittedOut := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")

		// yt-dlp returned nothing
		if strings.TrimSpace(splittedOut[0]) == "" {
			return "", hit
		}

		// If it's only one video, use the usual downloader
		if len(splittedOut) < 1 {
			cacheVideo[filename] = &[]*tele.Video{downloadSingleVideo(link, filename)}
		} else {
			var count int

			cacheVideo[filename] = &[]*tele.Video{}

			for _, l := range splittedOut {
				var data ytdlp
				err := json.Unmarshal([]byte(l), &data)
				if err == nil {
					// We need to proxy the requests for telegram when downloading
					*cacheVideo[filename] = append(*cacheVideo[filename], downloadSingleVideo(data.Url, filename, "--playlist-items", strconv.Itoa(count+1)))
					count++
				}
			}
		}
	}

	return filename, hit
}

func downloadSingleVideo(link string, filename string, arguments ...string) *tele.Video {
	// Starts yt-dlp with the arguments to select the best audio
	arguments = append(arguments, "-f", "bestvideo+bestaudio", "-f", "mp4", "-q", "-a", "-", "--geo-bypass", "-o", "-")
	ytDlp := exec.Command("yt-dlp", arguments...)
	ytDlp.Stdin = strings.NewReader(link)
	out, _ := ytDlp.StdoutPipe()
	_ = ytDlp.Start()

	go func() {
		err := ytDlp.Wait()
		if err != nil {
			lit.Error(err.Error())
		}
	}()

	return &tele.Video{File: tele.FromReader(out), FileName: filename, MIME: "video/mp4"}
}

func downloadAudio(link string) (string, bool) {
	hit := true

	filename := idGen(link) + ".mp3"

	if _, ok := cacheAudio[filename]; !ok {
		// Starts yt-dlp with the arguments to select the best audio
		ytDlp := exec.Command("yt-dlp", "-f", "bestaudio", "-f", "mp3", "-q", "-a", "-", "--geo-bypass", "-o", "-")
		ytDlp.Stdin = strings.NewReader(link)
		out, _ := ytDlp.StdoutPipe()
		_ = ytDlp.Start()

		cacheAudio[filename] = &tele.Audio{File: tele.FromReader(out), FileName: filename, MIME: "audio/mp3"}

		go func() {
			err := ytDlp.Wait()
			if err != nil {
				lit.Error(err.Error())
			}
		}()

		hit = false
	}

	return filename, hit
}

func downloadTikTok(link string) (string, bool, Media) {
	filename := idGen(link)

	// Check cache
	if _, ok := cacheAlbum[filename]; ok {
		return filename, true, Album
	} else {
		filename += ".mp4"
		if _, ok = cacheVideo[filename]; ok {
			return filename, true, Video
		}
	}
	// Remove the last four characters from filename
	filename = filename[:len(filename)-4]

	u, err := url.ParseRequestURI(cfg.Downloader)
	if err != nil {
		return "", false, Video
	}

	u.Path = "/api"
	u.RawQuery = url.Values{"url": {link}}.Encode()

	// Post to downloader
	resp, err := http.Get(u.String())
	if err == nil && resp.StatusCode == http.StatusOK {
		var d Downloader
		_ = json.NewDecoder(resp.Body).Decode(&d)

		switch d.Type {
		case "video":
			filename += ".mp4"
			cacheVideo[filename] = &[]*tele.Video{{File: tele.FromURL(d.VideoData.NwmVideoUrlHQ), MIME: "video/mp4", FileName: filename}}
			return filename, false, Video
		case "image":
			album := make([]*tele.Photo, len(d.ImageData.NoWatermarkImageList))
			for i, img := range d.ImageData.NoWatermarkImageList {
				album[i] = &tele.Photo{File: tele.FromURL(img)}
			}
			cacheAlbum[filename] = &album
			return filename, false, Album
		}
	}

	return "", false, Video
}

// cleanURL removes tracking and other unnecessary parameters from a URL
func cleanURL(link string) string {
	u, _ := url.Parse(link)
	q := u.Query()

	q.Del("utm_source")
	q.Del("utm_medium")
	q.Del("utm_name")
	q.Del("feature")
	q.Del("igshid")
	q.Del("si")

	u.RawQuery = q.Encode()

	return u.String()
}

func selectAndDownload(text string) (string, bool, Media) {
	var (
		media         Media
		filename      string
		hit           bool
		useDownloader bool
	)

	useDownloader = strings.Contains(text, "tiktok.com")

	if useDownloader {
		// Use the downloader to get videos and albums from TikTok
		filename, hit, media = downloadTikTok(text)
	}

	if !useDownloader || filename == "" {
		filename, hit = downloadYtDlp(text)
		media = Video
	}

	return filename, hit, media
}
