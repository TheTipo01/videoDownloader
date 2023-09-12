package main

import (
	"encoding/json"
	"github.com/bwmarrin/lit"
	tele "gopkg.in/telebot.v3"
)

const cacheTable = "CREATE TABLE IF NOT EXISTS `cache` ( `filename` VARCHAR(15) NOT NULL, `video` TEXT NOT NULL, PRIMARY KEY (`filename`) )"

func execQuery(query string) {
	_, err := db.Exec(query)
	if err != nil {
		lit.Error("Error executing query, %s", err)
		return
	}
}

func load() {
	var (
		filename string
		bytes    []byte
	)

	rows, err := db.Query("SELECT filename, video FROM cache")
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
		cache[filename] = &video
	}
}

func save(video *tele.Video) {
	bytes, _ := json.Marshal(video)

	_, err := db.Exec("INSERT INTO cache (filename, video) VALUES (?, ?)", video.FileName, bytes)
	if err != nil {
		lit.Error("Error executing query, %s", err)
		return
	}
}
