package main

import "sort"

var conf tomlConfig
var statevar State
var rawSongData []byte

// Data for configuration file

type tomlConfig struct {
	ClientData       clientInfo  `toml:"raspberry_pis"`
	ClientParameters clientParms `toml:"client_parameters"`
	ServerParameters serverParms `toml:"server_parameters"`
}

type clientInfo struct {
	Clients string `toml:"clients"`
}

type clientParms struct {
	CheckUpWaitTime int `toml:"check_up_wait_time"`
	MaxSyncLag      int `toml:"max_sync_lag"`
}

type serverParms struct {
	MusicFolder         string `toml:"music_folder"`
	Port                int    `toml:"port"`
	TimeToNextSong      int    `toml:"time_to_next_song"`
	TimeToDisallowSkips int    `toml:"time_to_disallow_skips"`
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
	RemoteComputers  []MakeConfig
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
