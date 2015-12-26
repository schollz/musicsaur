---
title:  "Musicsaur setup full instructions [Python]"
date:   2015-12-18 18:18:00
description: Full step-by-step to get going with Python and musicsaur!
permalink: python/
---


## Before you begin

Make sure that you have Python2.7+ or Python3.4+ installed on your server computer. Then install the required files using

{% highlight bash %}
sudo pip install setuptools
git clone https://github.com/schollz/musicsaur.git
cd musicsaur
sudo python setup.py install
{% endhighlight %}

## Configure

Copy the configuration file to a new file

{% highlight bash %}
cp config-python.cfg config.cfg
{% endhighlight %}

and edit line #42 with the locations of your music folderse. There are other parameters to edit, if you feel so inclined, but you needn't just to get started. Then simply run with

## Run

{% highlight bash %}
python syncmusic.py
{% endhighlight %}

## Other additions

If you'd like to autostart some Raspberry Pis, click [here to set that up](/raspberry-pi/) and then edit the configuration file.
