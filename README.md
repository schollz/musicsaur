# μsic

Want to sync up your music on multiple computers? This accomplishes exactly that - using only a simple Python server and a browser. Simply run the Python server, and open up a browser on each computer you want to sync up - that's it!

This program is powered by [the excellent howler.js library from goldfire](https://github.com/goldfire/howler.js/). Essentially all the client computers [sync their clocks](http://www.mine-control.com/zack/timesync/timesync.html) and then try to start a song at the same time. Any dissimilarities between playback are also fixed, because the clients will automatically seek to the position of the server.

# Installation

Tested on Python2.7 and Python3.4. To install simply use

```bash
git clone https://github.com/schollz/sync-music-player.git
python setup.py install
```

## Usage

Start the server using

```bash
python syncmusic.py "C:/Your/folder/of/music"
```

which should print out something like

```bash
############################################################
# Starting server with 18 songs
# To use, open a browser to http://W.X.Y.Z:5000
############################################################
```

Your server is up and running! Now, for each computer that you want to play music from, just go and load up a browser to the url ```http://W.X.Y.Z:5000```. You will see the playlist and the music will automatically synchronize and start playing! 

### Some notes

- If you don't hear anything, the client is probably trying to synchronize. The browser automatically mutes when it goes out of sync to avoid the headache caused by mis-aligned audio. You can see synchronization progress in [your browser console](https://webmasters.stackexchange.com/questions/8525/how-to-open-the-javascript-console-in-different-browsers). 
- If you still dont' hear anything, and you're using Chrome browser on Android [you need change one of the flags in chrome to allow audio without gestures](http://android.stackexchange.com/questions/59134/enable-autoplay-html5-video-in-chrome). To do this, copy and paste this into your Chrome browser:

```bash
chrome://flags/#disable-gesture-requirement-for-media-playback
```

- If you want to play music from a Raspberry Pi, just type this command (works on headless):

```bash
xinit /usr/bin/midori -a http://W.X.Y.Z:5000
```

## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D

## Credits

* James Simpson and [Goldfire studios](http://goldfirestudios.com/blog/104/howler.js-Modern-Web-Audio-Javascript-Library) for their amazing [howler.js library](https://github.com/goldfire/howler.js/)
* Zach Simpson for [his paper on simple clock synchronization](http://www.mine-control.com/zack/timesync/timesync.html)


