package main

import (
	"encoding/json"
	"fmt"
	mp3 "github.com/badgerodon/mp3"
	id3 "github.com/mikkyang/id3-go"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var m map[string]Song

type SyncJSON struct {
	Current_song     string `json:"current_song"`
	Client_timestamp int    `json:"client_timestamp"`
	Next_song        int    `json:"next_song"`
	Is_playing       bool   `json:"is_playing"`
	Song_time        int    `json:"song_time"`
}

type Song struct {
	Title  string
	Artist string
	Album  string
	Year   string
	Path   string
	Length time.Duration
}

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

func getMp3Info(path string) {
	defer timeTrack(time.Now(), "getMp3Info")
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	duration, err := mp3.Length(file)
	if err != nil {
		panic(err)
	}

	mp3File, err := id3.Open(path)
	fmt.Println(mp3File.Title())
	fmt.Println(mp3File.Artist())
	fmt.Println(mp3File.Album())
	fmt.Println(mp3File.Year())

	fmt.Printf("Expected length to return %v, got %v", time.Duration(5067754912), duration)

	m = make(map[string]Song)
	m["Something"] = Song{
		Title:  "hi",
		Artist: "noone",
		Length: time.Duration(5067754912),
	}
	fmt.Printf("Length %v\n", m["Something"].Length)
	fmt.Printf("%T\n", time.Duration(10000))
}

func main() {
	getMp3Info("./static/test.mp3")

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
