# Import core packages
import time
import sys
import os
import fnmatch
import random
import sys
import socket
import shutil
from threading import Timer
from configparser import SafeConfigParser
import json

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
import tornado.ioloop
import tornado.web
from jinja2 import Template
from mutagen.mp3 import MP3
from mutagen.id3 import ID3

#####################
# GLOBAL VARIABLES
#####################

state = {}
songStartTimer = None
songStopTimer = None
parser = SafeConfigParser()
try:
    parser.read('config.cfg')
except:
    print("Problem parsing config.cfg - did you change something?")
    sys.exit(-1)    

#####################
# UTILITY FUNCTIONS
#####################


def getTime():
    """Returns time in milliseconds, similar to Javascript"""

    return int(time.time() * 1000)


def getPlaylistHtml():
    """Returns HTML for the playlist"""

    playlist_html = ""
    for i,path in enumerate(state['ordering']):
        cur_song = state['playlist'][path]['song_name']
        html = """<a type="controls" data-skip="%(i)s">%(song)s</a><br>"""
        if state['currently_playing_songname'] == cur_song:
            song = "<b>" + cur_song + "</b>"
        else:
            song = cur_song
        playlist_html += html % {'i': str(i), 'song': song}
    return playlist_html


def songStarts():
    """Runs when server decides a song starts"""
    logger = logging.getLogger('syncmusic:songStarts')
    logger.debug('Playing: ' + state['currently_playing_songname'])
    state['is_playing'] = True



def songOver():
    """Runs when server decides a song stops"""
    logger = logging.getLogger('syncmusic:songOver')
    logger.debug('Done playing: ' + state['currently_playing_songname'])
    state['is_playing'] = False
    nextSong(int(parser.get('server_parameters','time_to_next_song')), -1)


def nextSong(delay, skip):
    """ Main song control

    Sets global flags for which song is playing,
    loads the new song, and sets the timers for when
    the songs should start and end.
    """
    global songStartTimer
    global songStopTimer
    logger = logging.getLogger('syncmusic:nextSong')
    if (time.time() - state['last_activated'] > int(parser.get('server_parameters','time_to_disallow_skips')) 
            or not state['is_initialized']): 

        state['last_activated'] = time.time()

        if skip < 0:
            state['current_song'] += skip + 2
        else:
            state['current_song'] = skip
        if state['current_song'] >= len(state['ordering']):
            state['current_song'] = 0
        if state['current_song'] < 0:
            state['current_song'] = len(state['orering']) - 1
        current_song_path = state['ordering'][state['current_song']]

        shutil.copy(current_song_path,os.path.join(os.getcwd(),'static/sound.mp3'))
        state['currently_playing_songname'] = state['playlist'][current_song_path]['song_name']
        state['next_song_time'] = getTime() + delay * 1000
        logger.debug('next up: ' + state['currently_playing_songname'])
        logger.debug('time: ' + str(getTime()) +
                     ' and next: ' + str(state['next_song_time']))
        state['is_initialized'] = True
        if songStartTimer is not None:
            songStartTimer.cancel()
            songStopTimer.cancel()
        songStopTimer = Timer(
            float(
                state['next_song_time'] -
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
                state['next_song_time'] -
                getTime()) /
            1000.0,
            songOver,
            ())
        songStartTimer.start()

#################
# WEB ROUTES
#################

index_page = Template(open('templates/index.html','r').read())

class IndexPage(tornado.web.RequestHandler):
    """Main sign-in - /

    Server loads new song if not initialized, and
    then returns the rendered music control page
    """

    def get(self):
        if not state['is_initialized']:
            nextSong(int(parser.get('server_parameters','time_to_next_song')), state['current_song'])
        data = {}
        data['random_integer'] = random.randint(1000, 30000)
        data['playlist_html'] = getPlaylistHtml()
        data['is_playing'] = state['is_playing']
        data['message'] = 'Syncing...'
        data['is_index'] = True
        data['max_sync_lag'] = parser.get('client_parameters','max_sync_lag')
        data['check_up_wait_time'] = parser.get('client_parameters','check_up_wait_time')
        if state['debug']:
            index_page = Template(open('templates/index.html','r').read())
        self.write(index_page.render(data=data))

class SyncHandler(tornado.web.RequestHandler):
    """Syncing route - /sync

    POST request from main page with the client client_timestamp
    and current_song. Returns JSON containing the server client_timestamp
    and whether or not to load a new song.
    """

    def post(self):
        data = {}
        data['client_timestamp'] = int(self.get_argument('client_timestamp'))
        data['server_timestamp'] = getTime()
        data['next_song'] = state['next_song_time']
        if state['is_playing']:
            data['is_playing'] = (state['currently_playing_songname'] == self.get_argument('current_song'))
        else:
            data['is_playing'] = state['is_playing']
        data['current_song'] = state['currently_playing_songname']
        data['song_time'] = float(getTime() - state['next_song_time']) / 1000.0
        self.write(data)


class NextSongHandler(tornado.web.RequestHandler):
    """Syncing route - /sync

    POST request from main page with the client client_timestamp
    and current_song. Returns JSON containing the server client_timestamp
    and whether or not to load a new song.
    """

    def post(self):
        response = {'message': 'loading!'}
        skip = int(self.get_argument('skip'))
        nextSong(int(parser.get('server_parameters','time_to_next_song')), skip)
        self.write(response)

# Depreciated
# @app.route("/playing", methods=['GET', 'POST'])
# def playing():
#     """ Is playing route - /nextSong

#     POST request to tell server that client has started
#     playing a song. DEPRECATED.
#     """
#     global is_playing
#     response = {'message': 'loading!'}
#     if request.method == 'POST':
#         is_playing = True
#     return jsonify(response)

application = tornado.web.Application([
    (r"/", IndexPage),
    (r"/sync", SyncHandler),
    (r"/nextsong", NextSongHandler),
    (r'/static/(.*)', tornado.web.StaticFileHandler, {'path': './static'}),
])



##########
# MAIN
##########

if __name__ == "__main__":
    """Load the playlist, or let user know that one needs to be loaded"""
    # app.run(host='0.0.0.0')
    logger = logging.getLogger('syncmusic:nextSong')
    cwd = os.getcwd()

    state['playlist'] = {}
    state['ordering'] = []
    state['current_song'] = 0
    state['currently_playing_songname'] = ""
    if os.path.isfile('state.json'):
        state = json.load(open('state.json','r'))
    state['last_activated'] = 0
    state['next_song_time'] = 0
    state['is_playing'] = False
    state['is_initialized'] = False
    state['last_activated'] = 0
    state['debug'] = True

    folders_with_music = parser.get('server_parameters','music_folder').split(',')
    for folder_with_music in folders_with_music:
        # Load playlist
        folder_with_music = folder_with_music.strip()
        for root, dirnames, filenames in os.walk(folder_with_music):
            for filename in fnmatch.filter(filenames, '*.mp3'):
                path = os.path.join(root, filename)
                if path in state['ordering']:
                    continue
                state['ordering'].append(path)
                state['playlist'][path] = {}
                title = filename
                artist = 'unknown'
                album = 'unknown'
                try:
                    audiofile = ID3(path)
                    try:
                        title = audiofile['TIT2'].text[0]
                    except:
                        pass
                    try:
                        artist = audiofile['TPE1'].text[0]
                    except:
                        pass
                    try:
                        album = audiofile['TALB'].text[0]
                    except:
                        pass
                except Exception, e:
                    logger.error('Failed to get audio information for ' + path + ': '+ str(e))
                state['playlist'][path]['title'] = title
                state['playlist'][path]['artist'] = artist
                state['playlist'][path]['album'] = album
                state['playlist'][path]['song_name'] = album + ' - ' + title + ' by ' + artist
                os.chdir(cwd)

    if len(state['ordering']) == 0:
        print(
            "\n\nNo mp3s found.\nDid you specify a music folder in line 40 of config.cfg?\n\n")
        sys.exit(-1)

    print(json.dumps(state,indent=2))
    os.chdir(cwd)
    try:
        s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        s.connect(("gmail.com",80))
        ip_address = s.getsockname()[0]
        s.close()
    except:
        ip_address = "127.0.0.1"

    print("\n\n" +"#" * 60)
    print("# Starting server with " + str(len(state['ordering'])) + " songs")
    print("# To use, open a browser to http://" + ip_address + ":"+ parser.get('server_parameters','port') + "")
    print("# To stop server, use Ctl + C")
    print("#" * 60 +"\n\n")

    pi_clients = []
    if len(parser.get('raspberry_pis','clients')) > 2 and ip_address != '127.0.0.1':
        pi_clients = parser.get('raspberry_pis','clients').split(',')
        for pi_client in pi_clients:
            pi_client = pi_client.strip()
            try:
                os.system("ssh " + pi_client + " 'pkill -9 midori </dev/null > log 2>&1 &'")
                os.system("ssh " + pi_client + " 'xinit /usr/bin/midori -a " + ip_address + ":" + parser.get('server_parameters','port') + "/ </dev/null > log 2>&1 &'")
            except:
                print("Problem starting pi!")

    application.listen(int(parser.get('server_parameters','port')))
    try:
        tornado.ioloop.IOLoop.instance().start()
    except (KeyboardInterrupt, SystemExit):
        print('\nProgram shutting down...')
        print('Saving state...')
        with open('state.json','w') as f:
            f.write(json.dumps(state,indent=2))

        for pi_client in pi_clients:
            pi_client = pi_client.strip()
            print('Shutting down ' + pi_client + '...')
            try:
                os.system("ssh " + pi_client + " 'pkill -9 midori </dev/null > log 2>&1 &'")
            except:
                pass
        print('Stopping timers...')
        try:
            songStopTimer.cancel()
        except:
            pass
        try:
            songStartTimer.cancel()
        except:
            pass
        print('Exiting...')
        sys.exit(-1)
        raise
