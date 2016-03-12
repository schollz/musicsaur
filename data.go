package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"sort"
)

var conf tomlConfig
var statevar State
var rawSongData []byte

// Configuration file stuff

type tomlConfig struct {
	MusicFolders []string
	Autostart    map[string]ClientSSH
	Server       serverParamaters
	Client       clientParameters
}

type ClientSSH struct {
	User          string
	Server        string
	Key           string
	Port          string
	Password      string
	RemoteBrowser string
}

type clientParameters struct {
	CheckupWaitTime int
	MaxSyncLag      int
}

type serverParamaters struct {
	Port                int
	TimeToNextSong      int
	TimeToDisallowSkips int
	Random              bool
}

func setupConfiguration() {
	if _, err := toml.DecodeFile("config.cfg", &conf); err != nil {
		fmt.Println(err)
		return
	}
}

// Data for state

type State struct {
	SongMap          map[string]Song
	SongList         sort.StringSlice
	PathList         map[string]bool
	SongStartTime    int64
	IsPlaying        bool
	CurrentSong      string
	CurrentSongIndex int
	LastMuted        int64
	IsMuted          bool
	IPAddress        string
	Port             int
	IndexPage        string
	MusicExtension   string
}

// Data for Song

type SyncJSON struct {
	Current_song        string  `json:"current_song"`
	Client_timestamp    int64   `json:"client_timestamp"`
	Server_timestamp    int64   `json:"server_timestamp"`
	Is_playing          bool    `json:"is_playing"`
	Song_time           float64 `json:"song_time"`
	Song_start_time     int64   `json:"next_song"`
	Mute_button_clicked bool    `json:"mute_button_clicked"`
	Is_muted            bool    `json:"is_muted"`
}

type Song struct {
	Fullname string
	Title    string
	Artist   string
	Album    string
	Path     string
	Length   int64
}

type IndexData struct {
	PlaylistHTML    string
	RandomInteger   int64
	CheckupWaitTime int64
	MaxSyncLag      int64
}
