---
title:  "Musicsaur setup full instructions [Golang]"
date:   2015-12-18 18:18:00
description: Full step-by-step to get going with Golang and musicsaur!
permalink: golang/
---


## Before you begin

First download the latest instrcutions

{% highlight bash %}
git clone https://github.com/schollz/musicsaur
cd musicsaur
{% endhighlight %}

then install the required packages

{% highlight bash %}
go get github.com/mholt/caddy/caddy
go get github.com/tcolgate/mp3
go get github.com/bobertlo/go-id3/id3
go get github.com/BurntSushi/toml
go get gopkg.in/tylerb/graceful.v1
{% endhighlight %}


## Configure

Then copy the configuration file 

{% highlight bash %}
cp config-go.cfg config.cfg
{% endhighlight %}

and edit line #5 with your music folders. 

## Run

To run, simply use

{% highlight bash %}
go run *.go
{% endhighlight %}

and you are good to go!

## Other additions

If you'd like to autostart some Raspberry Pis, click [here to set that up](/raspberry-pi/) and then edit the configuration file.


