import time
import sys
import shutil
import os
import fnmatch
import random

import eyed3
from flask import *

app = Flask(__name__)
app.debug = True


playlist = []
current_song = -1
last_activated = 0
next_song_time = 0
is_playing = False
is_initialized = False
song_name = ""

def getTime():
    return int(time.time()*1000)


@app.route("/")
def index_html():
    data = {}
    data['random_integer'] = random.randint(1000,30000)
    if not is_initialized:
        nextSong(20)
    if is_playing or next_song_time - getTime() > 17000:
        data = {}
        if is_playing:
            data['message'] = 'A song is currently playing. You will join on the next song.'
        else:
            data['message'] = 'Waiting for other participants, please hold.'
        data['is_index'] = False
    else:
        data['message'] = 'Starting soon.'
        data['is_index'] = True    
    return render_template('index.html',data = data)

@app.route("/sync", methods=['GET', 'POST'])
def sync():
    #searchword = request.args.get('key', '')
    if request.method == 'POST':
        print(getTime())
        data = {}
        data['client_timestamp'] = int(request.form['client_timestamp'])
        data['server_timestamp'] = getTime()
        data['next_song'] = next_song_time
        data['is_playing'] = is_playing
        data['current_song'] = song_name
        return jsonify(data)


@app.route("/nextsong", methods=['GET', 'POST'])
def finished():
    response = {'message':'loading!'}
    if request.method == 'POST':
        if is_playing:
            nextSong(20)
    return jsonify(response)

@app.route("/playing", methods=['GET', 'POST'])
def playing():
    global is_playing
    response = {'message':'loading!'}
    if request.method == 'POST':
        is_playing = True
    return jsonify(response)

def nextSong(delay):
    global last_activated
    global current_song
    global next_song_time
    global is_playing
    global is_initialized
    global song_name
    if time.time() - last_activated > 5 or not is_initialized: # songs can only be skipped every 5 seconds
        if not is_initialized:
            for root, dirnames, filenames in os.walk('/home/zack/Music'):
                for filename in fnmatch.filter(filenames, '*.mp3'):
                    if 'Allen' in root or 'Allen' in filename:
                        playlist.append((root, filename))
            print(playlist)
        is_playing = False
        current_song += 1
        if current_song >= len(playlist):
            current_song = 0
        print(current_song)
        last_activated = time.time()
        cwd = os.getcwd()
        print(playlist[current_song][0])
        os.chdir(playlist[current_song][0])
        cmd = 'scp ' + playlist[current_song][1].replace(' ','\ ') + ' phi@server8.duckdns.org:/www/data/sound.mp3'
        cmd = 'cp ' + playlist[current_song][1].replace(' ','\ ') + ' ' + cwd + '/static/sound.mp3'
        print(cmd)
        os.system(cmd)
        audiofile = eyed3.load(playlist[current_song][1])
        os.chdir(cwd)
        song_name = audiofile.tag.album + ' - ' + audiofile.tag.title + ' by ' + audiofile.tag.artist 
        next_song_time = getTime() + delay*1000
        print ('next up: ' + song_name)
        print ('time: ' + str(getTime()) + ' and next: ' + str(next_song_time))
        is_initialized = True

if __name__ == "__main__":
    # Load playlist
    app.run(host='10.190.76.50')

    from tornado.wsgi import WSGIContainer
    from tornado.httpserver import HTTPServer
    from tornado.ioloop import IOLoop
    http_server = HTTPServer(WSGIContainer(app))
    http_server.listen(5000)
    IOLoop.instance().start()

