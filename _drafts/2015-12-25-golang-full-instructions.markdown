---
title:  "Musicsaur setup full instructions [Golang]"
date:   2015-12-18 18:18:00
description: Full step-by-step to get going with Golang and MusicSAUR!
---


# Before you begin

You can use these instructions to get going on any kind of computer. For the sake of instruction though, I will assume you have two Raspberry Pis, pi1 and pi2, and that you want to have music playing on both of them with pi1 being the server. 

Make sure that you have Python2.7+ or Python3.4+ installed on your server computer (pi1 in this instructions).

# Setup

## Install midori (Required for Raspberry Pis only) 

First log into both of your `pi1` using

{% highlight bash %}
ssh pi1@pi1.ip.address
{% endhighlight %}

and install the midori webbrowser using

{% highlight bash %}
sudo apt-get install midori
{% endhighlight %}

Repeat for `pi2`.

## Install musicsaur on server

Now log back into `pi1` and run the following to install musicsaur:

{% highlight bash %}
sudo pip install setuptools
git clone https://github.com/schollz/musicsaur.git
cd musicsaur
sudo python setup.py install
{% endhighlight %}

## Configure

Now edit `config.cfg` with `vim` or `nano` and edit the line #42 to 

    music_folder = /location/of/your/music

and then edit line #16 to

{% highlight bash %}
clients = pi1@127.0.0.1,pi2@pi2.ip.address
{% endhighlight %}

where `pi2.ip.address` is the local IP addresses of your other pi. Now, still logged in to the pi1, you need to transfer a ssh key so you don't have to type a password in to start up music on these pis. To do this just type the following:

{% highlight bash %}
ssh-keygen
{% endhighlight %}

Press enter at each prompt. Now, from the same pi1 enter the following and type the requested password for each pi when it prompts:

{% highlight bash %}
ssh-copy-id pi1@127.0.0.1
ssh-copy-id pi2@pi2.ip.address
{% endhighlight %}

## Run it!

Now that you have uploaded your keys to each of your pi from the server pi you can have automatic SSH access. Now, still logged into to the server pi1, just type the following

{% highlight bash %}
python syncmusic.py
{% endhighlight %}

It will show an address that you can goto on your phone/computer to control and play the music. It will also automatically start up the Raspberry Pis browsers and they should start playing music.

## Other notes

If you want to run it and log off, you can use nohup mode:

{% highlight bash %}
nohup python syncmusic.py &
{% endhighlight %}

which can be killed with

{% highlight bash %}
pkill -9 python
{% endhighlight %}

To check whether the Raspberry Pis are working, simply ssh back into them and use

{% highlight bash %}
tail -f ~/log
{% endhighlight %}

to monitor what is going on.



[github]: https://github.com/schollz/musicsaur