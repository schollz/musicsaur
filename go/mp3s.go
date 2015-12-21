package main

import (
	"fmt"
	id3 "github.com/mikkyang/id3-go"
	mp3 "github.com/tcolgate/mp3"
	"os"
	"path/filepath"
	"time"
)

func loadMp3s(path string) {
	defer timeTrack(time.Now(), "loadMp3s")
	searchDir, _ := filepath.Abs("/home/zack/Music/Damien Rice/")

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
			fmt.Println(file)
			s := getMp3Info(file)
			songMap[s.Artist+" - "+s.Album+" - "+s.Title] = s
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

	mp3File, err := id3.Open(path)

	return Song{
		Title:  mp3File.Title(),
		Artist: mp3File.Artist(),
		Album:  mp3File.Album(),
		Path:   path,
		Length: getMp3Length(path),
	}
}

func getMp3Length(path string) (totalTime int64) {
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
