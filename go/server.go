package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/tylerb/graceful.v1"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
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
	//curTime = time.Since(time.Date(2015, 6, 1, 12, 0, 0, 0, time.UTC)).Nanoseconds() / 1000000
	return
}

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
		if statevar.CurrentSong != statevar.SongMap[k].Fullname {
			playlist_html += "<a type='controls' data-skip='" + strconv.Itoa(i) + "'>" + statevar.SongMap[k].Fullname + "</a><br>\n"
		} else {
			playlist_html += "<a type='controls' data-skip='" + strconv.Itoa(i) + "'><b>" + statevar.SongMap[k].Fullname + "</b></a><br>\n"

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
		data := SyncJSON{
			Current_song:        statevar.CurrentSong,
			Client_timestamp:    int64(client_timestamp),
			Server_timestamp:    getTime(),
			Is_playing:          statevar.IsPlaying,
			Song_time:           getPlaybackPositionInSeconds(),
			Song_start_time:     statevar.SongStartTime,
			Mute_button_clicked: mute_button_clicked,
			Is_muted:            is_muted,
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
	defer timeTrack(time.Now(), r.RemoteAddr+" /sync")
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
		statevar.CurrentSongIndex += song_index + 2
	} else {
		statevar.CurrentSongIndex = song_index
	}
	song := statevar.SongList[statevar.CurrentSongIndex]
	rawSongData, _ = ioutil.ReadFile(statevar.SongMap[song].Path)
	statevar.CurrentSong = statevar.SongMap[song].Fullname
	statevar.SongStartTime = getTime() + 11000
	statevar.IsPlaying = false
	b, _ := json.Marshal(statevar)
	ioutil.WriteFile("state.json", b, 0644)
	go songControl(statevar.SongStartTime-getTime()-3000, false, "3", statevar.SongMap[song].Fullname, false)
	go songControl(statevar.SongStartTime-getTime()-2000, false, "2", statevar.SongMap[song].Fullname, false)
	go songControl(statevar.SongStartTime-getTime()-1000, false, "1", statevar.SongMap[song].Fullname, false)
	go songControl(statevar.SongStartTime-getTime(), true, "Playing "+statevar.SongMap[song].Fullname, statevar.SongMap[song].Fullname, false)
	go songControl(statevar.SongStartTime-getTime()+statevar.SongMap[song].Length, false, "Stopping "+statevar.SongMap[song].Fullname, statevar.SongMap[song].Fullname, true)
}

func cleanup() {
	fmt.Println("cleanup")
}

func main() {

	piFlag := flag.String("pis", "", "\"pi@url1,pi@url2\"")
	portFlag := flag.String("port", "5000", "port to host on")
	libraryFlag := flag.String("folder", "./", "Folder to find the mp3s")
	flag.Parse()
	fmt.Println("piFlag:", *piFlag)
	fmt.Println("piFlag:", len(*piFlag))
	fmt.Println("portFlag:", *portFlag)
	fmt.Println("libraryFlag:", *libraryFlag)

	// // Load configuration parameters
	// if _, err := toml.DecodeFile("./config.cfg", &conf); err != nil {
	// 	// handle error
	// }
	// fmt.Printf("%v", conf)

	// Load state
	if _, err := os.Stat("state.json"); err == nil {
		dat, err := ioutil.ReadFile("state.json")
		if err != nil {
			panic(err)
		}
		json.Unmarshal(dat, &statevar)
		fmt.Println("\n*******")
		fmt.Println(statevar.CurrentSong)
		fmt.Println("*******\n")
		statevar.IsPlaying = false
		statevar.SongList = []string{}
		statevar.LastMuted = 0
		statevar.IsMuted = false
	} else {
		statevar = State{
			SongMap:          make(map[string]Song),
			SongList:         []string{},
			PathList:         make(map[string]bool),
			SongStartTime:    0,
			IsPlaying:        false,
			CurrentSong:      "None",
			CurrentSongIndex: 0,
			LastMuted:        0,
			IsMuted:          false,
		}
	}

	// Load Mp3s
	loadMp3s(*libraryFlag)

	// Load song list
	for k, _ := range statevar.SongMap {
		statevar.SongList = append(statevar.SongList, k)
	}
	statevar.SongList.Sort()

	skipTrack(statevar.CurrentSongIndex)
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//defer timeTrack(time.Now(), r.RemoteAddr+" /")
		html_response := index_html
		html_response = strings.Replace(html_response, "{{ data['random_integer'] }}", strconv.Itoa(rand.Intn(10000)), -1)
		html_response = strings.Replace(html_response, "{{ data['check_up_wait_time'] }}", strconv.Itoa(1700), -1)
		html_response = strings.Replace(html_response, "{{ data['max_sync_lag'] }}", strconv.Itoa(50), -1)
		html_response = strings.Replace(html_response, "{{ data['message'] }}", "Syncing...", -1)
		html_response = strings.Replace(html_response, "{{ data['playlist_html'] | safe }}", getPlaylistHTML(), -1)
		fmt.Fprintf(w, html_response)
	})

	mux.HandleFunc("/sound.mp3", func(w http.ResponseWriter, r *http.Request) {
		defer timeTrack(time.Now(), r.RemoteAddr+" /sound.mp3")
		w.Header().Set("Content-Type", "audio/mpeg")
		w.Write([]byte(rawSongData))
	})
	mux.HandleFunc("/static/howler.js", func(w http.ResponseWriter, r *http.Request) {
		//defer timeTrack(time.Now(), r.RemoteAddr+" /howler.js")
		w.Header().Set("Content-Type", "text/javascript")
		w.Write([]byte(howler_js))
	})
	mux.HandleFunc("/static/jquery.min.js", func(w http.ResponseWriter, r *http.Request) {
		//defer timeTrack(time.Now(), r.RemoteAddr+" /math.js")
		w.Header().Set("Content-Type", "text/javascript")
		w.Write([]byte(jquery_js))
	})
	mux.HandleFunc("/static/math.min.js", func(w http.ResponseWriter, r *http.Request) {
		//defer timeTrack(time.Now(), r.RemoteAddr+" /jquery.js")
		w.Header().Set("Content-Type", "text/javascript")
		w.Write([]byte(math_js))
	})
	mux.HandleFunc("/static/skeleton.css", func(w http.ResponseWriter, r *http.Request) {
		//defer timeTrack(time.Now(), r.RemoteAddr+" /skeleton.css")
		w.Header().Set("Content-Type", "text/css")
		w.Write([]byte(skeleton_css))
	})
	mux.HandleFunc("/static/normalize.css", func(w http.ResponseWriter, r *http.Request) {
		//defer timeTrack(time.Now(), r.RemoteAddr+" /normalize.css")
		w.Header().Set("Content-Type", "text/css")
		w.Write([]byte(normalize_css))
	})
	mux.HandleFunc("/sync", SyncRequest)
	mux.HandleFunc("/nextsong", NextSongRequest)
	//http.ListenAndServe(":5000", nil)

	ip, err := externalIP()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("\n\n######################################################################")
	fmt.Printf("# Starting server with %d songs\n", len(statevar.SongList))
	fmt.Println("# To use, open a browser to http://" + ip + ":" + *portFlag)
	fmt.Println("# To stop server, use Ctl + C")
	fmt.Println("######################################################################\n\n")

	graceful.Run(":"+*portFlag, 10*time.Second, mux)
}
