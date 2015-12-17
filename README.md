# sync-music-player

Allows you to play music from your library, in sync, with various other computers. This is the simplest form of syncing. Multiple clients request songs from the server, to which the server provides, along with a specified time to play the song. All the clients use the server to [sync their clocks](http://www.mine-control.com/zack/timesync/timesync.html) and then they all try to play the songs at the same time (works pretty well).

## Install

```bash
git clone https://github.com/schollz/sync-music-player.git
python setup.py install
```

## Server

Start the server using

```bash
python syncmusic.py "C:/Your/folder/of/music"
```

Then goto a browser and type in your ```localhost:5000``` to see the playlist.

## Client

There are two clients. The simple one is just to open a web browser and goto [localhsot:5000](http://localhost:5000/), or use whatever your host address is (probabaly 192.168.X.Y).  Note: If you are using Android, you won't be able to "autoplay" the music [unless you change one of the flags in chrome](http://android.stackexchange.com/questions/59134/enable-autoplay-html5-video-in-chrome). To do this, copy and paste this into your Chrome browser

```bash
chrome://flags/#disable-gesture-requirement-for-media-playback
```
and enable it.

To run on Raspberry Pi headless, use

```bash
xinit /usr/bin/midori -a http://ADDRESS:5000/
```

## Limitations

The main limitation is the upload of the music file, which may be a bottleneck if you are using it over the internet.

## Audio on Raspberry Pi



- [Maybe useful article about reducing crackle](https://dbader.org/blog/crackle-free-audio-on-the-raspberry-pi-with-mpd-and-pulseaudio#update1)
- [How to play MP3s with SOX](http://superuser.com/questions/421153/how-to-add-a-mp3-handler-to-sox/421168)

## Todo

- ~~Add Next, Previous buttons~~
- ~~Show playlist on every screen~~
- Mute the audio if it loses sync
