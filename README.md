![Logo](http://rpiai.com/musicsaur/musicsaur1.png)

# MusicSAUR (formerly the un-googleable "μsic")
## Music Synchronization And Uniform Relaying
[![Version 1.2](https://img.shields.io/badge/version-1.2-brightgreen.svg)]()
[![Join the chat at https://gitter.im/schollz/musicsaur](https://badges.gitter.im/schollz/music.svg)](https://gitter.im/schollz/musicsaur?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

---

**Update 12/19/2015: Version 1.2 released!**

**Major features include:**

- Improved server performance
- State variables for reloading position in playlist and faster boot-up
- Bug fixes
- Fewer dependencies

---

Want to sync up your music on multiple computers? This accomplishes exactly that - using only a simple Python server and a browser. Simply run the Python server, and open up a browser on each computer you want to sync up - that's it!

This program is powered by [the excellent howler.js library from goldfire](https://github.com/goldfire/howler.js/). Essentially all the client computers [sync their clocks](http://www.mine-control.com/zack/timesync/timesync.html) and then try to start a song at the same time. Any dissimilarities between playback are also fixed, because the clients will automatically seek to the position of the server.

# Installation

Tested on Python2.7 and Python3.4. Requires ```setuptools```. To install simply use

```bash
pip install setuptools
git clone https://github.com/schollz/sync-music-player.git
python setup.py install
```

## Usage

Open up the configuration file, ```config.cfg```, and edit line 42 by inputting your music folder, e.g.:

```bash
music_folder = /my/music/folder
```

Now start the server using

```bash
python syncmusic.py
```

which should print out something like

```bash
############################################################
# Starting server with 346 songs
# To use, open a browser to http://W.X.Y.Z:5000
# To stop server, use Ctl + C
############################################################
```

Your server is up and running! Now, for each computer that you want to play music from, just go and load up a browser to the url ```http://W.X.Y.Z:5000```. You will see the playlist and the music will automatically synchronize and start playing! 

### Auto-start Raspberry Pis!

If you'd like the server to automatically start up the Raspberry Pis, its easy to do. First, from your server computer (which can also be a Raspberry Pi) copy your ssh-key using the following command:

```bash
ssh-copy-id pi@YOUR_PIS_IP
```

Now sign-in to your Pi (you shouldn't have to use a password now) and install ```midori```:

```bash
ssh pi@YOUR_PIS_IP
sudo apt-get install midori
```

Finally, open up the ```config.cfg``` and change line 16 to include the address of the computer:

```bash
clients = pi@YOUR_PIS_IP
```

If you have more then one, just seperate them by commas.

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

OR

```bash
xinit /usr/bin/luakit -u http://W.X.Y.Z:5000
```


## Screenshot

![Screenshot](http://rpiai.com/musicsaur/screenshot2.png)

## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D

# History

- 12/19/2015 (morning): Version 1.2 Release
- 12/18/2015 (evening): Version 1.1 Release
- 12/18/2015 (morning): Version 1.0 Release


## Credits

* James Simpson and [Goldfire studios](http://goldfirestudios.com/blog/104/howler.js-Modern-Web-Audio-Javascript-Library) for their amazing [howler.js library](https://github.com/goldfire/howler.js/)
* Zach Simpson for [his paper on simple clock synchronization](http://www.mine-control.com/zack/timesync/timesync.html)
* Everyone on the [/r/raspberry_pi](https://www.reddit.com/r/raspberry_pi/comments/3xc8kq/simple_python_script_to_allow_multiple_raspberry/) and [/r/python](https://www.reddit.com/r/Python/comments/3xc8mj/simple_python_script_to_allow_multiple_computers/) threads for great feature requests and bug reports!
* [ClkerFreeVectorImages](https://pixabay.com/en/users/ClkerFreeVectorImages-3736/) and [OpenClipartVectors](https://pixabay.com/en/users/OpenClipartVectors-30363/) for the Public Domain vectors


