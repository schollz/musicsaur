package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func getTime() (curTime int64) {
	curTime = time.Now().UnixNano() / 1000000
	return
}

func songControl(millisecondWait int64, is_playing bool, text string, song string, start_next bool) {
	time.Sleep(time.Duration(millisecondWait) * time.Millisecond)
	if song == currentSong {
		log.Printf(song + " " + text)
		isPlaying = is_playing
		if start_next == true {
			skipTrack(1)
		}
	}
}

func getPlaylistHTML() (playlist_html string) {
	playlist_html = ""
	for i, k := range songList {
		if currentSong != k {
			playlist_html += "<a type='controls' data-skip='" + strconv.Itoa(i) + "'>" + k + "</a><br>\n"
		} else {
			playlist_html += "<a type='controls' data-skip='" + strconv.Itoa(i) + "'><b>" + k + "</b></a><br>\n"

		}
	}
	return
}

func getPlaybackPositionInSeconds() float64 {
	position := float64(getTime()-songStartTime) / 1000.0
	if isPlaying == true && position > 0 {
		return position
	} else {
		return 0.0
	}
}

func SyncRequest(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		//current_song := r.FormValue("current_song")
		client_timestamp, _ := strconv.Atoi(r.FormValue("client_timestamp"))
		data := SyncJSON{
			Current_song:     currentSong,
			Client_timestamp: int64(client_timestamp),
			Server_timestamp: getTime(),
			Is_playing:       isPlaying,
			Song_time:        getPlaybackPositionInSeconds(),
			Song_start_time:  songStartTime,
		}
		b, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(b))
	}
}

func NextSongRequest(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		skip, _ := strconv.Atoi(r.FormValue("skip"))
		skipTrack(skip)
		data := SyncJSON{
			Current_song:     "None",
			Client_timestamp: 0,
			Server_timestamp: 0,
			Is_playing:       false,
			Song_time:        0,
			Song_start_time:  0,
		}
		b, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.Write([]byte(b))
	}
}

func skipTrack(song_index int) {
	if song_index < 0 {
		currentSongIndex += song_index + 2
	} else {
		currentSongIndex = song_index
	}
	song := songList[currentSongIndex]
	rawSongData, _ = ioutil.ReadFile(songMap[song].Path)
	currentSong = song
	songStartTime = getTime() + 11000
	go songControl(songStartTime-getTime()-3000, false, "3", song, false)
	go songControl(songStartTime-getTime()-2000, false, "2", song, false)
	go songControl(songStartTime-getTime()-1000, false, "1", song, false)
	go songControl(songStartTime-getTime(), true, "Playing "+song, song, false)
	go songControl(songStartTime-getTime()+songMap[song].Length, false, "Stopping "+song, song, true)
}

func main() {
	currentSong = "None"
	currentSongIndex = 0
	isPlaying = false

	songMap = make(map[string]Song)
	loadMp3s("/home/zack/Music/Damien Rice/")
	fmt.Printf("%v\n", songMap)

	songList = []string{}
	for k, _ := range songMap {
		songList = append(songList, k)
	}
	songList.Sort()
	fmt.Println(songList)

	skipTrack(0)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer timeTrack(time.Now(), r.RemoteAddr+" /")
		html_response := index_html
		html_response = strings.Replace(html_response, "{{ .RandomInteger }}", strconv.Itoa(rand.Intn(10000)), -1)
		html_response = strings.Replace(html_response, "{{ .CheckupWaitTime }}", strconv.Itoa(1700), -1)
		html_response = strings.Replace(html_response, "{{ .MaxSyncLag }}", strconv.Itoa(50), -1)
		html_response = strings.Replace(html_response, "{{ .PlaylistHTML }}", getPlaylistHTML(), -1)
		html_response = strings.Replace(html_response, "{{ .Message }}", "hi", -1)
		fmt.Fprintf(w, html_response)
	})

	http.HandleFunc("/sound.mp3", func(w http.ResponseWriter, r *http.Request) {
		defer timeTrack(time.Now(), r.RemoteAddr+" /sound.mp3")
		log.Println("TESTING")
		w.Header().Set("Content-Type", "audio/mpeg")
		w.Write([]byte(rawSongData))
		// file, err := os.Open("./sound.mp3")
		// if err != nil {
		// 	panic(err)
		// }
		// defer file.Close()
		// http.ServeContent(w, r, "/static/test.mp3", time.Now(), file)

	})
	http.HandleFunc("/howler.js", func(w http.ResponseWriter, r *http.Request) {
		defer timeTrack(time.Now(), r.RemoteAddr+" /howler.js")
		w.Header().Set("Content-Type", "text/javascript")
		w.Write([]byte(howler_js))
	})
	http.HandleFunc("/math.js", func(w http.ResponseWriter, r *http.Request) {
		defer timeTrack(time.Now(), r.RemoteAddr+" /math.js")
		w.Header().Set("Content-Type", "text/javascript")
		w.Write([]byte(jquery_js))
	})
	http.HandleFunc("/jquery.js", func(w http.ResponseWriter, r *http.Request) {
		defer timeTrack(time.Now(), r.RemoteAddr+" /jquery.js")
		w.Header().Set("Content-Type", "text/javascript")
		w.Write([]byte(math_js))
	})
	http.HandleFunc("/skeleton.css", func(w http.ResponseWriter, r *http.Request) {
		defer timeTrack(time.Now(), r.RemoteAddr+" /skeleton.css")
		w.Header().Set("Content-Type", "text/css")
		w.Write([]byte(skeleton_css))
	})
	http.HandleFunc("/normalize.css", func(w http.ResponseWriter, r *http.Request) {
		defer timeTrack(time.Now(), r.RemoteAddr+" /normalize.css")
		w.Header().Set("Content-Type", "text/css")
		w.Write([]byte(normalize_css))
	})
	http.HandleFunc("/sync", SyncRequest)

	panic(http.ListenAndServe(":5000", nil))
}
