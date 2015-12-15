# sync-music-player

Allows you to play music from your library, in sync, with various other computers.

## Install

```bash
python setup.py install
```

## Configure and run

Load your own library by editing the YAML:

```yaml
folders:
	C:/some/folder/with/music
	C:/some/other/folder/with/music

filters:
	Allen Toussaint
	Rolling Stones
```

Then run the program using

```bash
python syncmusic.py
```

Then goto a browser and type in your ```yourip:5000``` to see the syncing.

## Limitations

The main limitation is the upload of the music file, which may be a bottleneck if you are using it over the internet.