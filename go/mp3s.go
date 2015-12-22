package main

import (
	"fmt"
	id3 "github.com/bobertlo/go-id3/id3"
	mp3 "github.com/tcolgate/mp3"
	"os"
	"path/filepath"
	"time"
)

func loadMp3s(path string) {
	fmt.Println("Looking in " + path)
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
				fmt.Println(file)
				s := getMp3Info(file)
				statevar.PathList[file] = true
				statevar.SongMap[s.Fullname] = s
			}
		}
	}
}

func getMp3Info(path string) Song {
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
		panic(err)
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
	}
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
