---
title:  "Getting started with Raspberry Pis"
date:   2015-12-24 16:18:00
description: How to use Raspberry Pi as a server/client
permalink: raspberry-pi/
---


It's very easy to get going on a Raspberry Pi as a server or a client. I'll go through the steps to set up a client and then the server. You can use a server computer also as a client computer, of course!

# Client

On the client computer install the ```midori``` browser using

{% highlight bash %}
sudo apt-get install midori
{% endhighlight %}

Now get the IP address using

{% highlight bash %}
ifconfig
{% endhighlight %}

and writing down the IP adress that starts 192.168... That's it!

# Server

On the server computer its easiest if you use the SSH keys. To generate an SSH key just use

{% highlight bash %}
ssh-keygen
{% endhighlight %}

Press enter at each prompt. Now transfer the keys to each other device (including this one, if the server is also going to be a client) using

{% highlight bash %}
ssh-copy-id piname@ipadress
{% endhighlight %}

Now your ready to install the server software. You can either use the Python version, the Golang version, or simply download the binary to run directly! Whatever version you use, you need to edit the configuration file with the information about the Raspberry Pi clients.

# Monitoring your client

To monitor whether the client is working, and to see error messages, simply log into the client and type

{% highlight bash %}
tail -f ~/log
{% endhighlight %}

to monitor what is going on.

