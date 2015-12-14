import time
import sys
import numpy

import requests
from flask import Flask, request, jsonify

app = Flask(__name__)
app.debug = True


next_song = time.time() + 15


@app.route("/")
def hello():
    return "Hello World!"

@app.route("/sync", methods=['GET', 'POST'])
def sync():
    #searchword = request.args.get('key', '')
    if request.method == 'POST':
        print(time.time())
        data = {}
        data['client_timestamp'] = request.form['client_timestamp']
        data['server_timestamp'] = time.time()
        data['next_song'] = next_song
        return jsonify(data)

def client():
    import pygame
    pygame.init()
    pygame.mixer.music.load("test.wav")
    correct_time_delta = []
    next_trigger = 0
    for i in range(5):
        print(time.time())
        r = requests.post('http://192.168.1.11:5000/sync',data = {'client_timestamp':time.time()})
        data = r.json()
        latency = time.time() - float(data['client_timestamp'])
        half_latency = latency / 2.0
        time_delta = time.time() - float(data['server_timestamp'])
        next_trigger = float(data['next_song'])
        correct_time_delta.append(time_delta + half_latency)
        time.sleep(1)
    print(correct_time_delta)
    start_time= time.time() + 10
    print('next trigger: ' + str(next_trigger))
    print('cur time' + str(time.time() - numpy.mean(correct_time_delta)))
    while True:
        if time.time() - numpy.mean(correct_time_delta) > next_trigger:
            print('done!')
            print(time.time() - numpy.mean(correct_time_delta))
            break
    pygame.mixer.music.play()
    while pygame.mixer.music.get_busy() == True:
        continue
    print('done playing music ' + str(time.time() - numpy.mean(correct_time_delta)))
    time.sleep(20)



if __name__ == "__main__":
    if len(sys.argv)<2:
        print "python timesync.py client/server IP"
        sys.exit(1)
    if sys.argv[1]=='server':
        print('next song:' + str(next_song))
        app.run(host='192.168.1.11')
    else:
        client()
