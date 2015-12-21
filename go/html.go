package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func test(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		current_song := r.FormValue("current_song")
		client_timestamp, _ := strconv.Atoi(r.FormValue("client_timestamp"))
		data := SyncJSON{
			Current_song:     "New song",
			Client_timestamp: client_timestamp,
			Next_song:        client_timestamp + 5000,
			Is_playing:       true,
			Song_time:        client_timestamp - 5000,
		}
		b, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%v\n", data)
		fmt.Printf("%v\n", b)
		fmt.Printf("%v %v", current_song, client_timestamp-2)
		rw.Write([]byte(b))
	}
}

func main() {

	songMap = make(map[string]Song)
	loadMp3s("./static/test.mp3")
	songList := make([]string, 0, len(songMap))
	for k := range songMap {
		songList = append(songList, k)
	}
	fmt.Printf("%v\n", songMap)
	fmt.Printf("%v\n", songList)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer timeTrack(time.Now(), r.RemoteAddr+" /")
		fmt.Fprintf(w, index_html)
	})
	http.HandleFunc("/static/test.mp3", func(w http.ResponseWriter, r *http.Request) {
		defer timeTrack(time.Now(), r.RemoteAddr+" /static/test.mp3")
		file, err := os.Open("./static/test.mp3")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		http.ServeContent(w, r, "/static/test.mp3", time.Now(), file)
	})
	http.HandleFunc("/howler.js", func(w http.ResponseWriter, r *http.Request) {
		defer timeTrack(time.Now(), r.RemoteAddr+" /howler.js")
		fmt.Fprintf(w, howler_js)
	})
	http.HandleFunc("/math.js", func(w http.ResponseWriter, r *http.Request) {
		defer timeTrack(time.Now(), r.RemoteAddr+" /math.js")
		fmt.Fprintf(w, jquery_js)
	})
	http.HandleFunc("/jquery.js", func(w http.ResponseWriter, r *http.Request) {
		defer timeTrack(time.Now(), r.RemoteAddr+" /jquery.js")
		fmt.Fprintf(w, math_js)
	})
	http.HandleFunc("/sync", test)

	panic(http.ListenAndServe(":17901", nil))
}
