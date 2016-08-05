package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"

	"strconv"
	"strings"
	"time"

	"github.com/mholt/caddy"
	"github.com/toqueteos/webbrowser"
	"gopkg.in/tylerb/graceful.v1"
)

const (
	appName    = "musicsaur"
	appVersion = "1.4.2"
)

func cleanup() {
	fmt.Println("cleanup")
}

func loadCaddyfile() (caddy.Input, error) {

	// Caddyfile in cwd
	contents := `IPADDRESS:PORT1 {
	proxy / IPADDRESS:PORT2 {
	except /static
}
	header  / Access-Control-Allow-Origin "*"
	tls off
	root ./
	gzip
	log ./caddy.log
}`
	contents = strings.Replace(contents, "IPADDRESS", statevar.IPAddress, -1)
	contents = strings.Replace(contents, "PORT1", strconv.Itoa(statevar.Port), -1)
	contents = strings.Replace(contents, "PORT2", strconv.Itoa(statevar.Port+1), -1)
	return caddy.CaddyfileInput{
		Contents:       []byte(contents),
		Filepath:       "Caddyfile",
		ServerTypeName: "http",
	}, nil
}

// RuntimeArgs contains all runtime
// arguments available
var RuntimeArgs struct {
	Port      string
	ServerCRT string
	ServerKey string
}

func main() {
	flag.StringVar(&RuntimeArgs.Port, "p", "8033", "port to bind")
	flag.StringVar(&RuntimeArgs.ServerCRT, "crt", "", "location of ssl crt")
	flag.StringVar(&RuntimeArgs.ServerKey, "key", "", "location of ssl key")
	flag.CommandLine.Usage = func() {
		fmt.Println(`musicsaur (version ` + appVersion + `): A Websocket Wiki and Kind Of A List Application
run this to start the server and then visit localhost at the port you specify
(see parameters).
Example: 'musicsaur -p 5000 127.0.0.1'
Options:`)
		flag.CommandLine.PrintDefaults()
	}

	dat, _ := ioutil.ReadFile("./static/logo.txt")
	fmt.Println(string(dat))
	fmt.Printf("\n\nversion %s\n", appVersion)

	setupConfiguration()

	// Load state
	if _, err := os.Stat("state.json"); err == nil {
		dat, err := ioutil.ReadFile("state.json")
		if err != nil {
			panic(err)
		}
		json.Unmarshal(dat, &statevar)
		statevar.IsPlaying = false
		statevar.SongList = []string{}
		statevar.LastMuted = 0
		statevar.IsMuted = false
		statevar.IndexPage = ""
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
			IndexPage:        "",
		}
	}

	// Parse flags
	flag.Parse()
	if flag.Arg(0) == "" {
		statevar.IPAddress = GetLocalIP()
	} else {
		statevar.IPAddress = flag.Arg(0)
	}

	port, _ := strconv.Atoi(RuntimeArgs.Port)
	statevar.Port = port

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//defer timeTrack(time.Now(), r.RemoteAddr+" /")
		// Load index page
		index_contents, err := ioutil.ReadFile("./templates/index.html")
		if err != nil {
			panic(err)
		}
		statevar.IndexPage = string(index_contents)
		html_response := statevar.IndexPage
		html_response = strings.Replace(html_response, "{{ data['random_integer'] }}", strconv.Itoa(rand.Intn(10000)), -1)
		html_response = strings.Replace(html_response, "{{ data['check_up_wait_time'] }}", strconv.Itoa(conf.Client.CheckupWaitTime), -1)
		html_response = strings.Replace(html_response, "{{ data['max_sync_lag'] }}", strconv.Itoa(conf.Client.MaxSyncLag), -1)
		html_response = strings.Replace(html_response, "{{ data['message'] }}", "Syncing...", -1)
		html_response = strings.Replace(html_response, "{{ data['playlist_html'] | safe }}", getPlaylistHTML(), -1)
		html_response = strings.Replace(html_response, "{{ data['sound_extension'] }}", statevar.MusicExtension, -1)
		fmt.Fprintf(w, html_response)
	})

	mux.HandleFunc("/sound.webm", func(w http.ResponseWriter, r *http.Request) {
		defer timeTrack(time.Now(), r.RemoteAddr+" /sound.mp3")
		w.Header().Set("Content-Type", "audio/mpeg")
		w.Write([]byte(rawSongData))
	})
	mux.HandleFunc("/sound.mp3", func(w http.ResponseWriter, r *http.Request) {
		defer timeTrack(time.Now(), r.RemoteAddr+" /sound.mp3")
		w.Header().Set("Content-Type", "audio/mpeg")
		w.Write([]byte(rawSongData))
	})
	mux.HandleFunc("/sync", SyncRequest)
	mux.HandleFunc("/nextsong", NextSongRequest)
	//http.ListenAndServe(":5000", nil)

	for _, k := range conf.Autostart {
		fmt.Println(k)
		response, err := runSSHCommand(k, "pkill -9 midori </dev/null > log 2>&1 &")
		fmt.Println(response)
		fmt.Println(err)
	}
	for _, k := range conf.Autostart {
		fmt.Println("Running autostart...")
		fmt.Println(k)
		cmd := "xinit /usr/bin/midori -a http://" + statevar.IPAddress + ":" + strconv.Itoa(statevar.Port) + "/ </dev/null > log 2>&1 &"
		fmt.Println(cmd)
		response, err := runSSHCommand(k, cmd)
		fmt.Println(response)
		fmt.Println(err)
	}

	fmt.Println("Open a browser to http://" + statevar.IPAddress + ":" + strconv.Itoa(statevar.Port))
	fmt.Println("To stop server, use Ctl + C\n\n")

	log.Printf("Starting server with %d songs\n", len(statevar.SongList))

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

	// start server

	skipTrack(statevar.CurrentSongIndex)
	webbrowser.Open("http://" + statevar.IPAddress + ":" + strconv.Itoa(statevar.Port))
	go graceful.Run(":"+strconv.Itoa(statevar.Port+1), 10*time.Second, mux)

	caddy.AppName = "Sprocket"
	caddy.AppVersion = "1.2.3"
	caddy.RegisterCaddyfileLoader("flag", caddy.LoaderFunc(confLoader))
	caddy.SetDefaultCaddyfileLoader("default", caddy.LoaderFunc(defaultLoader))
	// Get Caddyfile input
	caddyfile, err := caddy.LoadCaddyfile("http")
	if err != nil {
		log.Fatal(err)
	}

	// Start your engines
	fmt.Printf("%v", caddyfile)
	instance, err := caddy.Start(caddyfile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v", caddyfile)
	// Twiddle your thumbs
	instance.Wait()
}

// confLoader loads the Caddyfile using the -conf flag.
func confLoader(serverType string) (caddy.Input, error) {

	contents, err := ioutil.ReadFile("Caddyfile")
	if err != nil {
		return nil, err
	}
	return caddy.CaddyfileInput{
		Contents:       contents,
		Filepath:       "Caddyfile",
		ServerTypeName: serverType,
	}, nil
}

// defaultLoader loads the Caddyfile from the current working directory.
func defaultLoader(serverType string) (caddy.Input, error) {
	contents, err := ioutil.ReadFile(caddy.DefaultConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return caddy.CaddyfileInput{
		Contents:       contents,
		Filepath:       caddy.DefaultConfigFile,
		ServerTypeName: serverType,
	}, nil
}
