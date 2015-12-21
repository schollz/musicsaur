package main

import "sort"

var songMap map[string]Song
var songList sort.StringSlice
var songStartTime int64
var isPlaying bool
var currentSong string
var currentSongIndex int
var rawSongData []byte

type SyncJSON struct {
	Current_song     string  `json:"current_song"`
	Client_timestamp int64   `json:"client_timestamp"`
	Server_timestamp int64   `json:"server_timestamp"`
	Is_playing       bool    `json:"is_playing"`
	Song_time        float64 `json:"song_time"`
	Song_start_time  int64   `json:"song_start_time"`
}

type Song struct {
	Fullname string
	Title    string
	Artist   string
	Album    string
	Path     string
	Length   int64
}

type IndexData struct {
	PlaylistHTML    string
	RandomInteger   int64
	CheckupWaitTime int64
	MaxSyncLag      int64
}

var index_html2 = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">

    <title>MusicSAUR</title>
    <meta content="width=device-width, initial-scale=1" name="viewport">
    <meta content="no-cache" http-equiv="Cache-control">
    <meta content="-1" http-equiv="Expires">
    <script src="/math.js" type="text/javascript">
    </script>
    <script src="/jquery.js" type="text/javascript">
    </script>
    <script src="/howler.js" type="text/javascript">
    </script>
</head>
<body>
<audio controls preload="auto" src="./sound.mp3" id="sound" type="audio/mpeg">
  Your browser does not support the audio tag.
</audio>
</body>
</html>
`

var index_html = `
		<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">

    <title>MusicSAUR</title>
    <meta content="width=device-width, initial-scale=1" name="viewport">
    <meta content="no-cache" http-equiv="Cache-control">
    <meta content="-1" http-equiv="Expires">
    <script src="/math.js" type="text/javascript">
    </script>
    <script src="/jquery.js" type="text/javascript">
    </script>
    <script src="/howler.js" type="text/javascript">
    </script>
    <!--
<audio preload="auto" src="/sound.mp3?{{ .RandomInteger }}" id="sound" type="audio/mpeg">
  Your browser does not support the audio tag.
</audio>-->
    <link href="/normalize.css" rel="stylesheet">
    <link href="/skeleton.css" rel="stylesheet">
    <style>
    a { cursor: pointer; }

    .u-pull-right {
  float: left; }
    </style>
    <script>

var sound = new Howl({
  src: ['/sound.mp3?{{ .RandomInteger }}'],
  preload: true
});

    </script>
</head>

<body>
<script>

var time = Date.now || function() {
return +new Date.getTime();
}

// CONSTANTS
var CHECK_UP_WAIT_TIME = {{ .CheckupWaitTime }};
var CHECK_UP_ITERATION = 1;
var check_up_counter = 0;
var MAX_SYNC_LAG = {{ .MaxSyncLag }};

// GLOBALS
var lagTimes = [];
var tryWait = 0;
var computeTimes = [];
var correct_time_delta = [];
var correct_latency = [];
var next_trigger = time() + 1000000;
var true_time_delta = 0;
var true_server_time_delta = 0;
var sound_activated = false;
var seconds_left = 0;
var current_song = "None"
var current_song_name = "None"
var secondTimeout3 = setTimeout(function() {
console.log('3 seconds left')
}, 100000);
var secondTimeout2 = setTimeout(function() {
console.log('2 seconds left')
}, 100000);
var secondTimeout1 = setTimeout(function() {
console.log('1 seconds left')
}, 100000);
var secondTimeout0 = setTimeout(function() {
console.log('0 seconds left')
}, 100000);
var mainInterval = 0;
var runningDiff = 0;



function makeRequests(callback) {

for (var i = 0; i < 23; i++) {

  setTimeout(function postRequest() {

// Send the data using post
    var posting = $.post('/sync', {
        'client_timestamp': time(),
        'current_song': current_song
    });

    // Put the results in a div
    posting.done(function(data) {
    	console.log(data);
    	console.log(data['current_song'])
        var timeNow = time();
        current_song = data['current_song']
        latency = timeNow - data['client_timestamp']
        half_latency = latency / 2.0
        time_delta = timeNow - data['server_timestamp']
        next_trigger = data['song_start_time']

        correct_time_delta.push(time_delta + half_latency);
        correct_latency.push(half_latency);
        if (correct_time_delta.length==23) {
          console.log('correct_time_delta');
          console.log(correct_time_delta);
          var mean = math.mean(correct_time_delta);
          var median = math.median(correct_time_delta);
          var std = math.std(correct_time_delta);
          var sum = 0
          var num = 0
          for (var j = 0; j < correct_time_delta.length; j++) {
              if (correct_time_delta[j]<median+std) {
                  sum = sum + correct_time_delta[j];
                  num = num + 1;
              }
          }
          true_time_delta = sum / num;

          var mean = math.mean(correct_latency);
          var median = math.median(correct_latency);
          var std = math.std(correct_latency);
          var sum = 0
          var num = 0
          for (var j = 0; j < correct_latency.length; j++) {
              if (correct_latency[j]<median+std) {
                  sum = sum + correct_latency[j];
                  num = num + 1;
              }
          }
          true_server_time_delta = sum / num;

          clearTimeout(secondTimeout3);
          secondTimeout3 = setTimeout(function() {
              console.log('3 seconds left');
              $("div.info1").text('Playing in 3...');
          }, next_trigger - (time() - true_time_delta) - 3000);
          clearTimeout(secondTimeout2);
          secondTimeout2 = setTimeout(function() {
              console.log('2 seconds left');
              $("div.info1").text('Playing in 2...');
          }, next_trigger - (time() - true_time_delta) - 2000);
          clearTimeout(secondTimeout1);
          secondTimeout1 = setTimeout(function() {
              console.log('1 seconds left');
              $("div.info1").text('Playing in 1...');
          }, next_trigger - (time() - true_time_delta) - 1000);
          clearTimeout(secondTimeout0);
          secondTimeout0 = setTimeout(function() {
              console.log('playing song');
              current_song_name = current_song.split(":");
              current_song_name = current_song_name[current_song_name.length-1];
              $("div.info1").html('Loading <b>' + current_song_name + '</b>...');
              mainInterval = setInterval(function(){
                checkIfSkipped();
              }, CHECK_UP_WAIT_TIME);
              sound.play();
              if (data['is_playing']==true) {
                  sound.seek(data['song_time'])
              }
              // var posting = $.post('/playing', {
              // 'message': 'im playing a song'
              // });

              // // Put the results in a div
              // posting.done(function(data) {

              // });
          }, next_trigger - (time() - true_time_delta));

        }
    });

  }, i*180 );
    
}
}


function checkIfSkipped() {

    
    // Send the data using post
    var posting = $.post('/sync', {
        'client_timestamp': time(),
        'current_song': current_song
    });

    // Put the results in a div
    posting.done(function(data) {
        var start = new Date().getTime();
        var time_delta2 = time()-(data['server_timestamp']+true_time_delta);
      check_up_counter = check_up_counter + 1;
      if (data['is_playing']==false) {
        console.log('reloading page');
        sound.unload()
        location.reload(true);
      } else if (check_up_counter %% CHECK_UP_ITERATION==0) {
        check_up_counter = 0;
        var mySongTime = sound.seek();
        if (typeof(mySongTime)=="object") {
          mySongTime = 0;
          console.log('Still loading...')
            $("div.info1").html('Loading <b>' + current_song_name + '</b>...');
        }

        if (mySongTime == 0) {
          sound.seek(data['song_time']+time_delta2/1000.0);
        } else {
          var diff = data['song_time']+time_delta2/1000.0 - mySongTime;
          if (Math.abs(diff) > MAX_SYNC_LAG/1000.0) {
            CHECK_UP_ITERATION = 1;
            sound.volume(0.0);
            runningDiff = runningDiff + diff;
            var serverSongTime = data['song_time']+time_delta2/1000.0;
            console.log('[' + Date.now() + '] ' + ': NOT in sync (>' + MAX_SYNC_LAG.toString() + ' ms)')
            console.log('Browser:  ' + mySongTime.toString() + '\nServer: ' + serverSongTime.toString() + '\nDiff: ' + (diff*1000).toString() + '\nMean half-latency: ' + true_server_time_delta.toString() +  '\nMeasured half-latency: ' + time_delta2.toString() + '\nrunningDiff: ' + (runningDiff*1000).toString() + '\nSeeking to: ' + (serverSongTime+runningDiff).toString());
            $("div.info1").html('Muted <b>' + current_song_name + '</b> (out of sync)');
            if (diff<-1000000) {
              console.log('pausing')
              sound.pause()
              clearTimeout(secondTimeout3);
              clearTimeout(mainInterval);
              secondTimeout3 = setTimeout(function() {
                  console.log('playing');
                  sound.play();
                  mainInterval = setInterval(function(){
                    checkIfSkipped();
                  }, CHECK_UP_WAIT_TIME);
              }, Math.abs(runningDiff)*1000);
            } else {
                console.log(JSON.stringify(data));
                sound.seek(serverSongTime+runningDiff);
            }
          } else {
            console.log('[' + Date.now() + '] ' + ': in sync (|' + (diff*1000).toString() + '|<' + MAX_SYNC_LAG.toString() + ' ms)')
            $("div.info1").html('Playing <b>' + current_song_name + '</b>');
            CHECK_UP_ITERATION = parseInt(30.0/(CHECK_UP_WAIT_TIME/1000.0)); // every 30 seconds
            tryWait = 0;
            check_up_counter = 0;
            sound.volume(1.0);
          } 
        }
      }
    });

}




$(document).ready(function(){
$('a[type=controls]').click(function() {
   var skip = $(this).data('skip');
   console.log(skip);
    $("div.info1").text('Changing song');
    var posting = $.post('/nextsong', {
        'message': 'next song please',
        'skip': skip
    });

    // Put the results in a div
    posting.done(function(data) {
        sound.unload()
        location.reload(true);
        console.log('reloading page')
    });

});


makeRequests();

});




</script>

    <div class="container">
        <!-- columns should be the immediate child of a .row -->


        <div class="row">
        </div>


        <div class="row">
        <span style="display:table;">

            <span style="vertical-align: middle; display: table-cell;">
            <h1  style="position:relative;bottom:0"><i>musicsaur</i><br><small style="font-size: 50%%;">&nbsp;version 1.2</small></h1>
        </span>
    </div>


        <div class="row">
            <div class="seven columns">
                <a class="button" data-skip="-3" type="controls">Previous</a> <a class="button" data-skip="-2" type="controls">Replay</a> <a class="button" data-skip="-1" type="controls">Next</a>
            </div>


            <div class="five columns">
                <div class="info1">
                    {{ .Message }}
                </div>
            </div>
        </div>


        <div class="row">
            <div class="two columns">
            </div>


            <div class="ten columns">
                <div class="info2">
                    {{ .PlaylistHTML }}
                </div>
            </div>
        </div>
    </div>
</body>
</html>

`
