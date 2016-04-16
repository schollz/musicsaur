package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func songControl(millisecondWait int64, is_playing bool, text string, song string, start_next bool) {
	time.Sleep(time.Duration(millisecondWait) * time.Millisecond)
	if song == statevar.CurrentSong {
		log.Printf(song + " " + text)
		statevar.IsPlaying = is_playing
		if start_next == true {
			skipTrack(-1)
		}
	}
}

func getPlaylistHTML() (playlist_html string) {
	playlist_html = ""
	for i, k := range statevar.SongList {
		name := statevar.SongMap[k].Fullname
		names := strings.Split(name, "/")
		showName := names[len(names)-1]
		name = showName
		names = strings.Split(name, "\\")
		showName = names[len(names)-1]
		if statevar.CurrentSong != statevar.SongMap[k].Fullname {
			playlist_html += "<a type='controls' data-skip='" + strconv.Itoa(i) + "'>" + showName + "</a><br>\n"
		} else {
			playlist_html += "<a type='controls' data-skip='" + strconv.Itoa(i) + "'><b>" + showName + "</b></a><br>\n"

		}
	}
	return
}

func getPlaybackPositionInSeconds() float64 {
	position := float64(getTime()-statevar.SongStartTime) / 1000.0
	if statevar.IsPlaying == true && position > 0 {
		return position
	} else {
		return 0.0
	}
}

func SyncRequest(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// defer timeTrack(time.Now(), r.RemoteAddr+" /sync")
		//current_song := r.FormValue("current_song")
		client_timestamp_str := r.FormValue("client_timestamp")
		client_timestamp, _ := strconv.ParseUint(client_timestamp_str, 10, 64)
		is_muted, _ := strconv.ParseBool(r.FormValue("is_muted"))
		mute_button_clicked, _ := strconv.ParseBool(r.FormValue("mute_button_clicked"))
		if mute_button_clicked == true {
			statevar.LastMuted = getTime()
			statevar.IsMuted = is_muted
		}

		if getTime()-statevar.LastMuted < 3000 {
			mute_button_clicked = true
			is_muted = statevar.IsMuted
		}
		name := statevar.CurrentSong
		names := strings.Split(name, "/")
		showName := names[len(names)-1]
		name = showName
		names = strings.Split(name, "\\")
		showName = names[len(names)-1]
		data := SyncJSON{
			Current_song:        showName,
			Client_timestamp:    int64(client_timestamp),
			Server_timestamp:    getTime(),
			Is_playing:          statevar.IsPlaying,
			Song_time:           getPlaybackPositionInSeconds(),
			Song_start_time:     statevar.SongStartTime,
			Mute_button_clicked: mute_button_clicked,
			Is_muted:            statevar.IsMuted,
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
		defer timeTrack(time.Now(), r.RemoteAddr+" /nextsong")
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
		if song_index == -1 && conf.Server.Random {
			s1 := rand.NewSource(time.Now().UnixNano())
			r1 := rand.New(s1)
			statevar.CurrentSongIndex = r1.Intn(len(statevar.SongList))
		} else {
			statevar.CurrentSongIndex += song_index + 2
		}
	} else {
		statevar.CurrentSongIndex = song_index
	}
	if statevar.CurrentSongIndex >= len(statevar.SongList) {
		statevar.CurrentSongIndex = 0
	}
	log.Println(statevar.CurrentSongIndex, len(statevar.SongList))
	song := statevar.SongList[statevar.CurrentSongIndex]

	err := os.Remove("./static/sound.mp3")
	if err != nil {
		log.Println(err)
	}

	// To be served by Caddy
	CopyFile(statevar.SongMap[song].Path, "./static/sound.mp3")
	statevar.MusicExtension = "mp3"
	if conf.Server.Ffmpeg {

		start := time.Now()
		cmd := "ffmpeg"
		args := []string{"-i", "./static/sound.mp3", "-y", "-acodec", "pcm_u8", "-ar", "44100", "./static/sound.wav"}
		if err := exec.Command(cmd, args...).Run(); err != nil {
			// If unsuccessful, will defualt to sound.mp3
			log.Println("Error with mp3 -> wav", err)
		} else {
			elapsed := time.Since(start)
			log.Printf("mp3 -> wav done. (%s)\n", elapsed)
			start = time.Now()
			cmd = "ffmpeg"
			args = []string{"-i", "./static/sound.wav", "-y", "-dash", "1", "-c:a", "libopus", "-compression_level", "0", "-frame_duration", "60", "-application", "lowdelay", "-cutoff", "20000", "./static/sound.webm"}
			if err := exec.Command(cmd, args...).Run(); err != nil {
				// If unsuccessful, will defualt to sound.mp3
				log.Println("Error with wav -> webm", err)
			} else {
				// If successful use sound.webm
				elapsed = time.Since(start)
				log.Printf("wav -> webm done. (%s)\n", elapsed)
				statevar.MusicExtension = "webm"
			}
		}
	}

	rawSongData, _ = ioutil.ReadFile(statevar.SongMap[song].Path)

	statevar.CurrentSong = statevar.SongMap[song].Fullname
	statevar.SongStartTime = getTime() + 9000
	statevar.IsPlaying = false
	b, _ := json.Marshal(statevar)
	ioutil.WriteFile("state.json", b, 0644)
	go songControl(statevar.SongStartTime-getTime()-3000, false, "3", statevar.SongMap[song].Fullname, false)
	go songControl(statevar.SongStartTime-getTime()-2000, false, "2", statevar.SongMap[song].Fullname, false)
	go songControl(statevar.SongStartTime-getTime()-1000, false, "1", statevar.SongMap[song].Fullname, false)
	go songControl(statevar.SongStartTime-getTime(), true, "Playing "+statevar.SongMap[song].Fullname, statevar.SongMap[song].Fullname, false)
	go songControl(statevar.SongStartTime-getTime()+statevar.SongMap[song].Length, false, "Stopping "+statevar.SongMap[song].Fullname, statevar.SongMap[song].Fullname, true)
}
