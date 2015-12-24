package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/tylerb/graceful.v1"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	//"os/exec"
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
		defer timeTrack(time.Now(), r.RemoteAddr+" /sync")
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

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func skipTrack(song_index int) {
	if song_index < 0 {
		statevar.CurrentSongIndex += song_index + 2
	} else {
		statevar.CurrentSongIndex = song_index
	}
	song := statevar.SongList[statevar.CurrentSongIndex]

	err := os.Remove("./data/sound.mp3")
	if err != nil {
		fmt.Println(err)
	}

	// To be served by Caddy
	CopyFile(statevar.SongMap[song].Path, "./data/sound.mp3")

	// cmd := "cp"
	// args := []string{statevar.SongMap[song].Path, "/cygdrive/C/Users/ZNS/Desktop/Caddy/stuff/sound.mp3"}
	// if err := exec.Command(cmd, args...).Run(); err != nil {
	// 	fmt.Println(err)
	// 	fmt.Fprintln(os.Stderr, err)
	// 	os.Exit(1)
	// }
	// fmt.Println("Shrinking file...")
	// cmd := "ffmpeg"
	// err := os.Remove("sound.mp3")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// args := []string{"-i", statevar.SongMap[song].Path, "-codec:a", "libmp3lame", "-qscale:a", "8", "sound.mp3"}
	// if err := exec.Command(cmd, args...).Run(); err != nil {
	// 	fmt.Println(err)
	// 	fmt.Fprintln(os.Stderr, err)
	// 	os.Exit(1)
	// }
	// fmt.Println("Successfully shrunk file")
	// rawSongData, _ = ioutil.ReadFile("sound.mp3")

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

	setupConfiguration()

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
		if len(conf.MusicFolders) == 0 {
			executable := strings.Split(os.Args[0], "\\")
			executable_name := executable[len(executable)-1]
			fmt.Println("Run \"" + executable_name + " --help\" to learn how to add a folder of music")
			os.Exit(0)
		}
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
	if len(conf.MusicFolders) > 0 {
		for _, folder := range conf.MusicFolders {
			loadMp3s(folder)
		}
	}

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
		html_response = strings.Replace(html_response, "{{ data['check_up_wait_time'] }}", strconv.Itoa(conf.Client.CheckupWaitTime), -1)
		html_response = strings.Replace(html_response, "{{ data['max_sync_lag'] }}", strconv.Itoa(conf.Client.MaxSyncLag), -1)
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

	ip := GetLocalIP()

	port := strconv.Itoa(conf.Server.Port)

	fmt.Println("\n\n######################################################################")
	fmt.Printf("# Starting server with %d songs\n", len(statevar.SongList))
	fmt.Println("# To use, open a browser to http://" + ip + ":" + port)
	fmt.Println("# To stop server, use Ctl + C")
	fmt.Println("######################################################################\n\n")

	for _, k := range conf.Autostart {
		fmt.Println(k)
		response, err := runSSHCommand(k, "pkill -9 midori </dev/null > log 2>&1 &")
		fmt.Println(response)
		fmt.Println(err)
	}
	for _, k := range conf.Autostart {
		fmt.Println(k)
		cmd := "xinit /usr/bin/midori -a http://" + ip + ":" + port + "/ </dev/null > log 2>&1 &"
		fmt.Println(cmd)
		response, err := runSSHCommand(k, cmd)
		fmt.Println(response)
		fmt.Println(err)
	}

	graceful.Run(":"+port, 10*time.Second, mux)
}
