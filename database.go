package main

import (
	"encoding/json"
	"github.com/bwmarrin/lit"
	tele "gopkg.in/telebot.v3"
)

const videoTable = "CREATE TABLE IF NOT EXISTS `video` ( `filename` VARCHAR(15) NOT NULL, `video` TEXT NOT NULL, PRIMARY KEY (`filename`) )"
const albumTable = "CREATE TABLE IF NOT EXISTS `album` ( `filename` VARCHAR(11) NOT NULL, `album` TEXT NOT NULL, PRIMARY KEY (`filename`) )"
const audioTable = "CREATE TABLE IF NOT EXISTS `audio` ( `filename` VARCHAR(15) NOT NULL, `audio` TEXT NOT NULL, PRIMARY KEY (`filename`) )"

func execQuery(query ...string) {
	for _, q := range query {
		_, err := db.Exec(q)
		if err != nil {
			lit.Error("Error executing query, %s", err)
			return
		}
	}
}

func load() {
	var (
		filename string
		bytes    []byte
	)

	// Videos
	rows, err := db.Query("SELECT filename, video FROM video")
	if err != nil {
		lit.Error("Error executing query, %s", err)
		return
	}

	for rows.Next() {
		var video tele.Video

		err = rows.Scan(&filename, &bytes)
		if err != nil {
			lit.Error("Error scanning row, %s", err)
			continue
		}

		_ = json.Unmarshal(bytes, &video)
		cacheVideo[filename] = &video
	}

	// Photos
	rows, err = db.Query("SELECT filename, album FROM album")
	if err != nil {
		lit.Error("Error executing query, %s", err)
		return
	}

	for rows.Next() {
		var photos []*tele.Photo

		err = rows.Scan(&filename, &bytes)
		if err != nil {
			lit.Error("Error scanning row, %s", err)
			continue
		}

		err = json.Unmarshal(bytes, &photos)
		cacheAlbum[filename] = &photos
	}

	// Audio
	rows, err = db.Query("SELECT filename, audio FROM audio")
	if err != nil {
		lit.Error("Error executing query, %s", err)
		return
	}

	for rows.Next() {
		var audio tele.Audio

		err = rows.Scan(&filename, &bytes)
		if err != nil {
			lit.Error("Error scanning row, %s", err)
			continue
		}

		err = json.Unmarshal(bytes, &audio)
		cacheAudio[filename] = &audio
	}
}

func saveVideo(video *tele.Video) {
	bytes, _ := json.Marshal(video)

	_, err := db.Exec("INSERT INTO video (filename, video) VALUES (?, ?)", video.FileName, bytes)
	if err != nil {
		lit.Error("Error executing query, %s", err)
		return
	}
}

func saveAlbum(album *[]*tele.Photo, filename string) {
	bytes, _ := json.Marshal(album)

	_, err := db.Exec("INSERT INTO album (filename, album) VALUES (?, ?)", filename, bytes)
	if err != nil {
		lit.Error("Error executing query, %s", err)
		return
	}
}

func saveAudio(audio *tele.Audio) {
	bytes, _ := json.Marshal(audio)

	_, err := db.Exec("INSERT INTO audio (filename, audio) VALUES (?, ?)", audio.FileName, bytes)
	if err != nil {
		lit.Error("Error executing query, %s", err)
		return
	}
}
