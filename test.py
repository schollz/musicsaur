import time
import sys
import shutil
import os

from flask import *

app = Flask(__name__)
app.debug = True


playlist = ['short1.mp3','short2.mp3','long2.mp3','long3.mp3','short3.mp3']
current_song = -1
last_activated = 0
next_song_time = 0
is_playing = False
is_initialized = False

def getTime():
    return int(time.time()*1000)


@app.route("/")
def index_html():
    data = {}

    if not is_initialized:
        nextSong(20)
    if is_playing or next_song_time - getTime() > 7000:
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
        data['client_timestamp'] = request.form['client_timestamp']
        data['server_timestamp'] = getTime()
        data['next_song'] = next_song_time
        data['is_playing'] = is_playing
        data['current_song'] = playlist[current_song]
        return jsonify(data)


@app.route("/nextsong", methods=['GET', 'POST'])
def finished():
    response = {'message':'loading!'}
    if request.method == 'POST':
        nextSong(6)
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
    if time.time() - last_activated > 5 or not is_initialized: # songs can only be skipped every 5 seconds
        is_playing = False
        current_song += 1
        if current_song >= len(playlist):
            current_song = 0
        last_activated = time.time()
        # shutil.copy('./' + playlist[current_song],'./static/')
        # os.rename('./static/' + playlist[current_song],'./static/sound.mp3')
        # os.system('scp ' + playlist[current_song] + ' phi@192.168.1.11:/www/data/sound.mp3')
        next_song_time = getTime() + delay*1000
        print ('next up: ' + playlist[current_song])
        is_initialized = True

if __name__ == "__main__":
    app.run(host='10.190.76.50')

