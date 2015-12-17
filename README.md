# μsic

Want to sync up your music on multiple computers? This accomplishes exactly that - using only a simple Python server and a browser. Simply run the Python server, and open up a browser on each computer you want to sync up - that's it!

This program is powered by [the excellent howler.js library from goldfire](https://github.com/goldfire/howler.js/). Essentially all the client computers [sync their clocks](http://www.mine-control.com/zack/timesync/timesync.html) and then try to start a song at the same time. Any dissimilarities between playback are also fixed, because the clients will automatically seek to the position of the server.

## Install

```bash
git clone https://github.com/schollz/sync-music-player.git
python setup.py install
```

## Run

Start the server using

```bash
python syncmusic.py "C:/Your/folder/of/music"
```

Now, figure out your local server IP using ```ifconfig``` or similar. Then goto a browser and type in your ```http://LOCALSERVERIP:5000``` to see the playlist and hear the synced up music! Note: If you are using Android, you won't be able to hear the music [unless you change one of the flags in chrome to allow audio without gestures](http://android.stackexchange.com/questions/59134/enable-autoplay-html5-video-in-chrome). To do this, copy and paste this into your Chrome browser

```bash
chrome://flags/#disable-gesture-requirement-for-media-playback
```
and enable it. Also note, if you are using a Raspberry Pi, you can run the browser headless using the following command:

```bash
xinit /usr/bin/midori -a http://LOCALSERVERIP:5000/
```

## Todo

- ~~Add Next, Previous buttons~~
- ~~Show playlist on every screen~~
- ~~Mute the audio if it loses sync~~
