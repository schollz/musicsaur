package main

import (
	"encoding/json"
	"fmt"
	id3 "github.com/mikkyang/id3-go"
	mp3 "github.com/tcolgate/mp3"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var songMap map[string]Song

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
	Path   string
	Length int64
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

func main() {

	songMap = make(map[string]Song)
	loadMp3s("./static/test.mp3")
	fmt.Printf("%v\n", songMap)

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
