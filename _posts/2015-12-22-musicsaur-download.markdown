---
title:  "Musicsaur binary installation"
date:   2015-12-22 18:18:00
description: Instructions to just download and run.
permalink: download-binary/
---

<style>
a.cta3 {
	background: #5badf0;
	color: #fff;
	margin-left: 12px;
	padding: 8px 12px;
	font-size: 13px;
	/*font-weight: bold;*/
	line-height: 1.35;
	border-radius: 3px;
}
</style>

# Download

<table class="tg">
  <tr>
    <th class="tg-031e">OS</th>
    <th class="tg-yw4l">Download link</th>
    <th class="tg-yw4l">SHA1SUM</th>
  </tr>
  <tr>
    <td class="tg-031e">Windows (64-bit)</td>
    <td class="tg-yw4l"><a href="/assets/builds/musicsaur-1.3.0-windows-amd64.exe.zip">musicsaur-1.3.0-windows-amd64</a></td>
    <td class="tg-yw4l">0455797d4bc797eaf1c5cda1fd1cf4be1f9dd880</td>
  </tr>
  <tr>
    <td class="tg-yw4l">Mac OS X</td>
    <td class="tg-yw4l">
    <a href="/assets/builds/musicsaur-1.3.0-darwin-amd64.zip">musicsaur-1.3.0-darwin-amd64</a>
    </td>
    <td class="tg-yw4l">346177c8fb0329127894ea193626ac4a2fd09aad</td>
  </tr>
  <tr>
    <td class="tg-yw4l">Linux 64-bit</td>
    <td class="tg-yw4l">

    <a href="/assets/builds/musicsaur-1.3.0-linux-amd64.zip">musicsaur-1.3.0-linux-amd64</a></th>


    </td>
    <td class="tg-yw4l">be885cc4965fd8ff86ad0da0b08531653e779eb3</td>
  </tr>
  <tr>
    <td class="tg-yw4l">Raspberry Pi</td>
    <td class="tg-yw4l">

    <a href="/assets/builds/musicsaur-1.3.0-linux-arm.zip">musicsaur-1.3.0-linux-arm</a></th>


    </td>
    <td class="tg-yw4l">0b2831fc8f1b646775f1e0f3c83f8af204eb458d</td>
  </tr>
</table>


# Instructions

First download the zipped archive corresponding to your operating system (Windows, Mac, Linux, Raspberry Pi). Optionally, you can check the SHA1 sum of the archive to see if it matches correspondingly.

Then, unzip the archive into the folder of your choice. You will see two folders, a configuration file, and an executable. 

Now, you can open the configuration file, {% highlight bash %}config.cfg{% endhighlight %} using the editor of your choice. In line #5 of this file, you need to edit to include the path of your music. Make sure to include the *full path* and not the relative path. No other parameters need to be edited in this configuration file.

To run, simply double click on the executable, or run within a Terminal. Upon running, you should see something like this in the terminal window:

{% highlight bash %}
##########################################
# Starting server with 10 songs
# To use, open a browser to http://IP:PORT
# To stop server, use Ctl + C
##########################################
{% endhighlight %}

To see the music player now, simply open any browser and type in the URL {% highlight bash %}http://IP:PORT{% endhighlight %} to see the music controls.


# License

Copyright (c) 2015 Zack

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.