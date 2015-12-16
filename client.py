import time
import sys
import os
import random
import logging
import json
from threading import Timer,Thread,Event
import subprocess

class perpetualTimer():

	def __init__(self,t,hFunction):
		self.t=t
		self.hFunction = hFunction
		self.thread = Timer(self.t,self.handle_function)

	def handle_function(self):
		self.hFunction()
		self.thread = Timer(self.t,self.handle_function)
		self.thread.start()

	def start(self):
		self.thread.start()

	def cancel(self):
		self.thread.cancel()




root = logging.getLogger()
root.setLevel(logging.DEBUG)

ch = logging.StreamHandler(sys.stdout)
ch.setLevel(logging.DEBUG)
formatter = logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s')
ch.setFormatter(formatter)
root.addHandler(ch)

import requests
import numpy

url = ""

def getTime():
    return int(time.time()*1000)


def checkIfSkipped():
	logger = logging.getLogger('client:checkIfSkipped')
	print "Checking if skipped"
	request_url = url + '/sync'
	request_url = 'http://' + request_url.replace('//','/')

	data = {'client_timestamp':getTime()}
	try:
		r = requests.post(request_url,data=data)
		data = r.json()
		logger.info(json.dumps(data))
		if not data['is_playing']:
			os.system('pkill play')
			os.system('rm sound.mp3')
			os.system('wget http://' + url + '/static/sound.mp3')
			syncClocks()
	except:
		pass

def syncClocks():
	correct_time_delta = []
	request_url = url + '/'
	request_url = 'http://' + request_url.replace('//','/')
	r = requests.get(request_url)

	for i in range(5):
		request_url = url + '/sync'
		request_url = 'http://' + request_url.replace('//','/')

		data = {'client_timestamp':getTime()}
		r = requests.post(request_url,data=data)
		data = r.json()
		latency = getTime() - data['client_timestamp']
		half_latency = latency / 2.0
		time_delta = getTime() - data['server_timestamp']
		next_trigger = data['next_song']

		correct_time_delta.append(time_delta)
		time.sleep(0.15)

	mean_time_delta = numpy.mean(correct_time_delta)
	print(mean_time_delta)
	sleep_time = (next_trigger - (getTime() - mean_time_delta))/1000.0
	if sleep_time > 0:
		time.sleep(sleep_time-3)
		print('3')
		time.sleep(1)
		print('2')
		time.sleep(1)
		print('1')
		time.sleep(1)
		print('playing')
		subprocess.Popen(["play","sound.mp3"])
		request_url = url + '/playing'
		request_url = 'http://' + request_url.replace('//','/')
		data = {'client_timestamp':getTime()}
		r = requests.post(request_url,data=data)

if __name__ == "__main__":
	try:
		url=sys.argv[1]
	except:
		print('python client.py URL')
		sys.exit(-1)
	t = perpetualTimer(3,checkIfSkipped)
	t.start()
