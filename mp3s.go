package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	id3 "github.com/bobertlo/go-id3/id3"
	mp3 "github.com/tcolgate/mp3"
)

func loadMp3s(path string) {
	defer timeTrack(time.Now(), "loadMp3s")
	searchDir, _ := filepath.Abs(path)

	fileList := []string{}
	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	for _, file := range fileList {
		if filepath.Ext(file) == ".mp3" {
			if !statevar.PathList[file] {
				log.Println(file)
				s, err := getMp3Info(file)
				if err != nil {
					log.Println("Couldn't get ID3 for " + file + ", skipping...")
				} else {
					statevar.PathList[file] = true
					statevar.SongMap[s.Path] = s
				}
			}
		}
	}
}

func getMp3Info(path string) (Song, error) {
	defer timeTrack(time.Now(), "getMp3Info")
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// duration, err := mp3.Length(file)
	// if err != nil {
	// 	log.Println(err)
	// 	duration = time.Duration(10000)
	// }

	tags, err := id3.ReadFile(file)
	if err != nil {
		return Song{
			Fullname: "none",
			Title:    "none",
			Artist:   "none",
			Album:    "none",
			Path:     "none",
			Length:   0,
		}, err
	}
	title := tags["title"]
	artist := tags["artist"]
	album := tags["album"]
	fullname := artist + " - " + album + " - " + title
	if title == "" {
		title = path
		fullname = title
	}
	return Song{
		Fullname: fullname,
		Title:    title,
		Artist:   artist,
		Album:    album,
		Path:     path,
		Length:   getMp3Length(path),
	}, nil
}

func getMp3Length(path string) (totalTime int64) {
	// Returns length in milliseconds
	r, err := os.Open(path)
	if err != nil {
		//fmt.Println(err)
		return
	}

	d := mp3.NewDecoder(r)
	var f mp3.Frame
	totalTime = 0
	for {

		if err := d.Decode(&f); err != nil {
			//fmt.Println(err)
			return
		}

		totalTime += f.Duration().Nanoseconds() / 1000000
	}
}
