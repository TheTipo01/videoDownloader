package main

import (
	"crypto/sha256"
	"encoding/base32"
	"github.com/bwmarrin/lit"
	"net/url"
	"os"
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

func downloadVideo(link string) string {
	lit.Debug(link)

	fileName := idGen(link) + ".mp4"

	_, err := os.Stat(tempFolder + fileName)
	if err != nil {
		// Starts yt-dlp with the arguments to select the best audio
		ytDlp := exec.Command("yt-dlp", "-f", "bestvideo+bestaudio", "-f", "mp4", "-q", "-a", "-", "--geo-bypass", "--output", tempFolder+fileName)
		ytDlp.Stdin = strings.NewReader(link)
		out, err := ytDlp.CombinedOutput()
		if err != nil {
			lit.Error("Error while downloading video: %s", string(out))
		}
	}

	return fileName
}

func contains(str string, checkFor []string) bool {
	for _, s := range checkFor {
		if strings.Contains(str, s) {
			return true
		}
	}

	return false
}
