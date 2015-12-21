package main

var songMap map[string]Song

type SyncJSON struct {
	Current_song     string `json:"current_song"`
	Client_timestamp int    `json:"client_timestamp"`
	Next_song        int    `json:"next_song"`
	Is_playing       bool   `json:"is_playing"`
	Song_time        int    `json:"song_time"`
}

type Song struct {
	Title  string
	Artist string
	Album  string
	Path   string
	Length int64
}

var index_html = `
		<html>
		<head>
    <meta content="width=device-width, initial-scale=1" name="viewport">
    <meta content="no-cache" http-equiv="Cache-control">
    <meta content="-1" http-equiv="Expires">
		<script type="text/javascript" src="/howler.js"></script>
		<script type="text/javascript" src="/math.js"></script>
		<script type="text/javascript" src="/jquery.js"></script>
		</head>
		<body>
		<script>
		var sound = new Howl({
  src: ['/static/test.mp3'],
  autoplay: true,
  loop: true,
  volume: 0.5,
  onend: function() {
    console.log('Finished!');
  }
});
var time = Date.now || function() {
return +new Date.getTime();
}

	</script>
		<h1>Welcome Home!</h1>
		</body>
		</html>
		`
