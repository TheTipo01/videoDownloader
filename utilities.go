package main

import (
	"crypto/sha256"
	"encoding/base32"
	"github.com/bwmarrin/lit"
	tele "gopkg.in/telebot.v3"
	"net/url"
	"os/exec"
	"strings"
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

func checkAndDownload(link string) (string, bool) {
	hit := true
	lit.Debug(link)

	filename := idGen(link) + ".mp4"

	if _, ok := cache[filename]; !ok {
		// Starts yt-dlp with the arguments to select the best audio
		ytDlp := exec.Command("yt-dlp", "-f", "bestvideo+bestaudio", "-f", "mp4", "-q", "-a", "-", "--geo-bypass", "-o", "-")
		ytDlp.Stdin = strings.NewReader(link)
		out, _ := ytDlp.StdoutPipe()
		_ = ytDlp.Start()

		cache[filename] = &tele.Video{File: tele.FromReader(out), FileName: filename, MIME: "video/mp4"}

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
