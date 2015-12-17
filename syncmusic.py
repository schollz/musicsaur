# Import core packages
import time
import sys
import os
import fnmatch
import random
import sys
from threading import Timer

# Setup logging
import logging
root = logging.getLogger()
root.setLevel(logging.DEBUG)
ch = logging.StreamHandler(sys.stdout)
ch.setLevel(logging.DEBUG)
formatter = logging.Formatter(
    '%(asctime)s - %(name)s - %(levelname)s - %(message)s')
ch.setFormatter(formatter)
root.addHandler(ch)

# Import 3rd-part packages
import eyed3
from flask import *
from mutagen.mp3 import MP3


#####################
# GLOBAL VARIABLES
#####################

app = Flask(__name__)
app.debug = True

playlist = []
playlist_info = []
current_song = -1
last_activated = 0
next_song_time = 0
is_playing = False
is_initialized = False
song_name = ""
songStartTimer = None
songStopTimer = None
folder_with_music = ""

#####################
# UTILITY FUNCTIONS
#####################


def getTime():
    """Returns time in milliseconds, similar to Javascript"""

    return int(time.time() * 1000)


def getPlaylistHtml():
    """Returns HTML for the playlist"""

    playlist_html = ""
    for i in range(len(playlist_info)):
        html = """<a type="controls" data-skip="%(i)s">%(song)s</a><br>"""
        if playlist_info[i] == song_name:
            song = "<b>" + playlist_info[i] + "</b>"
        else:
            song = playlist_info[i]
        playlist_html += html % {'i': str(i), 'song': song}
    return playlist_html


def songStarts():
    """Runs when server decides a song starts"""
    logger = logging.getLogger('syncmusic:songStarts')
    logger.debug('Playing: ' + song_name)


def songOver():
    """Runs when server decides a song stops"""
    global is_playing
    logger = logging.getLogger('syncmusic:songOver')
    logger.debug('Done playing: ' + song_name)
    is_playing = False
    nextSong(6, -1)


def nextSong(delay, skip):
    """ Main song control

    Sets global flags for which song is playing,
    loads the new song, and sets the timers for when
    the songs should start and end.
    """
    global last_activated
    global current_song
    global next_song_time
    global is_playing
    global is_initialized
    global song_name
    global songStartTimer
    global songStopTimer
    logger = logging.getLogger('syncmusic:nextSong')
    if time.time() - last_activated > 3 or not is_initialized:  # songs can only be skipped every 5 seconds
        is_playing = False
        if skip < 0:
            current_song += skip + 2
        else:
            current_song = skip
        if current_song >= len(playlist):
            current_song = 0
        if current_song < 0:
            current_song = len(playlist) - 1

        last_activated = time.time()
        cwd = os.getcwd()
        os.chdir(playlist[current_song][0])
        cmd = 'cp ' + \
            playlist[current_song][1].replace(' ', '\ ') + ' ' + cwd + '/static/sound.mp3'
        logger.debug(cmd)
        os.system(cmd)
        os.chdir(cwd)
        song_name = playlist_info[current_song]
        next_song_time = getTime() + delay * 1000
        logger.debug('next up: ' + song_name)
        logger.debug('time: ' + str(getTime()) +
                     ' and next: ' + str(next_song_time))
        is_initialized = True
        if songStartTimer is not None:
            songStartTimer.cancel()
            songStopTimer.cancel()
        songStopTimer = Timer(
            float(
                next_song_time -
                getTime()) /
            1000.0,
            songStarts,
            ())
        songStopTimer.start()
        audio = MP3('./static/sound.mp3')
        logger.debug(audio.info.length)
        songStartTimer = Timer(
            2 +
            float(
                audio.info.length) +
            float(
                next_song_time -
                getTime()) /
            1000.0,
            songOver,
            ())
        songStartTimer.start()

#################
# WEB ROUTES
#################


@app.route("/")
def index_html():
    """Main sign-in - /

    Server loads new song if not initialized, and
    then returns the rendered music control page
    """

    if not is_initialized:
        nextSong(6, 0)
    data = {}
    data['random_integer'] = random.randint(1000, 30000)
    data['playlist_html'] = getPlaylistHtml()
    data['is_playing'] = is_playing
    data['message'] = 'Syncing...'
    data['is_index'] = True
    return render_template('index.html', data=data)


@app.route("/sync", methods=['GET', 'POST'])
def sync():
    """Syncing route - /sync

    POST request from main page with the client client_timestamp
    and current_song. Returns JSON containing the server client_timestamp
    and whether or not to load a new song.
    """
    #searchword = request.args.get('key', '')
    if request.method == 'POST':
        data = {}
        data['client_timestamp'] = int(request.form['client_timestamp'])
        data['server_timestamp'] = getTime()
        data['next_song'] = next_song_time
        if is_playing:
            data['is_playing'] = (song_name == request.form['current_song'])
        else:
            data['is_playing'] = is_playing
        data['current_song'] = song_name
        data['song_time'] = float(getTime() - next_song_time) / 1000.0
        return jsonify(data)


@app.route("/nextsong", methods=['GET', 'POST'])
def finished():
    """ Next song route - /nextSong

    POST request to start a new song
    """
    response = {'message': 'loading!'}
    if request.method == 'POST':
        skip = int(request.form['skip'])
        nextSong(6, skip)
    return jsonify(response)


@app.route("/playing", methods=['GET', 'POST'])
def playing():
    """ Is playing route - /nextSong

    POST request to tell server that client has started
    playing a song. DEPRECATED.
    """
    global is_playing
    response = {'message': 'loading!'}
    if request.method == 'POST':
        is_playing = True
    return jsonify(response)


##########
# MAIN
##########

if __name__ == "__main__":
    """Load the playlist, or let user know that one needs to be loaded"""
    # app.run(host='0.0.0.0')
    logger = logging.getLogger('syncmusic:nextSong')
    if len(sys.argv) > 1:
        folder_with_music = sys.argv[1]
        # Load playlist
        for root, dirnames, filenames in os.walk(folder_with_music):
            for filename in fnmatch.filter(filenames, '*.mp3'):
                playlist.append((root, filename))
                cwd = os.getcwd()
                os.chdir(root)
                audiofile = eyed3.load(filename)
                if audiofile.tag is None:
                    continue
                title = audiofile.tag.title
                if title is None:
                    title = 'unknown'
                artist = audiofile.tag.artist
                if artist is None:
                    artist = filename
                album = audiofile.tag.album
                if album is None:
                    album = ''
                song_name = album + ' - ' + title + ' by ' + artist
                playlist_info.append(song_name)
                os.chdir(cwd)
        if len(playlist) == 0:
            print('No music in ' + folder_with_music)
            sys.exit(-1)
    else:
        print(
            "Need to specify folder with music.\npython syncmusic.py '/folder/with/music'")
        sys.exit(-1)

    print("*" * 60)
    print("* Starting server with " + str(len(playlist)) + " songs")
    print("*" * 60)

    from tornado.wsgi import WSGIContainer
    from tornado.httpserver import HTTPServer
    from tornado.ioloop import IOLoop
    http_server = HTTPServer(WSGIContainer(app))
    http_server.listen(5000)
    IOLoop.instance().start()
