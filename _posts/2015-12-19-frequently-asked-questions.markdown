---
title:  "Frequently Asked Questions"
date:   2015-12-19 18:18:00
description: What is this program? What are the alternatives?
permalink: faq/
---

**What is musicsaur?**

> *musicsaur* is simply the simplest music synchronization program I know. All you do is run a program and open a browser on the computers that you want to synchronize together. Its fast, easy, and surprisingly good at what it does.

> The name, musicsaur, stands for Music Synchronization And Uniform Replay, which is the jargony way I like to describe musicsaur.

**Is it free?**

> Yes! It's free. Not only that, but its also [open-source](https://github.com/schollz/musicsaur), licensed by the MIT license.

**What are some alternatives to multi-room wifi playback?**

> - [Chromecast audio](http://www.androidcentral.com/chromecast-audio-can-now-play-same-song-every-room). They are [2 for $55](https://store.google.com/product/chromecast_audio).
> - Pianobar
> - [Volumio](https://volumio.org/)
> - [Pulseaudio](http://www.danplanet.com/blog/2014/11/26/multi-room-audio-with-multicast-rtp/)
> - [Shairport-Sync](https://github.com/mikebrady/shairport-sync)


**Can you add songs once the server starts?**

> Songs are only loaded when the server starts. That's a good thing to add though. I'd also like to have better organization of the playlist on the webpage (at least sorted by artist...)


**How come I don't hear anything?**

> If you don't hear anything, the client is probably trying to synchronize. The browser automatically mutes when it goes out of sync to avoid the headache caused by mis-aligned audio. You can see synchronization progress in [your browser console](https://webmasters.stackexchange.com/questions/8525/how-to-open-the-javascript-console-in-different-browsers)

**Does this work on a phone?**

> Yes, at least for Chrome on Android. To have it work, [you need change one of the flags in chrome to allow audio without gestures](http://android.stackexchange.com/questions/59134/enable-autoplay-html5-video-in-chrome). To do this, copy and paste this into your Chrome browser:

{% highlight bash %}
chrome://flags/#disable-gesture-requirement-for-media-playback
{% endhighlight %}



