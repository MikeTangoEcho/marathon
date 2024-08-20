# Marathon

Use the live chat to stream loop your videos.

## How it works

**Requirements**: Install FFmpeg
* https://www.ffmpeg.org/download.html

```sh
marathon start --twitch-stream-key [twitch stream key. ex: live_XXXX] --twitch-broadcaster-channel [twitch username, ex: justintv]
```

1. The bot connects to iRc and joins your channel, then wait for commands
2. On `!play lotr.playlist` => the bot will stream in real time your playlist with FFmpeg ffconcat.
3. On `!shutdown` => the bot will halt the stream and exits.

## What's Next ?

Streaming platform
* [x] Twitch
* [ ] YouTube

Broadcaster
* [x] FFmpeg *concat*
* [ ] vlc
* [ ] liquidsoap

## Usefull links

* https://help.twitch.tv/s/twitch-ingest-recommendation
* https://help.twitch.tv/s/article/broadcasting-guidelines
* https://ffmpeg.org/ffmpeg-all.html#toc-Preset-files
* https://ffmpeg.org/ffmpeg-protocols.html#udp
* https://trac.ffmpeg.org/wiki/Concatenate
* https://trac.ffmpeg.org/wiki/EncodingForStreamingSites
* https://trac.ffmpeg.org/wiki/StreamingGuide
* https://trac.ffmpeg.org/wiki/Encode/H.264

## Local Mock: RTMP server to HLS stream

**Start RTMP server that write an HLS stream**

https://ffmpeg.org/ffmpeg-formats.html#toc-hls-2

```sh
ffmpeg -f flv -listen 1 -i rtmp://127.0.0.1:1935/live/app  -hls_time 2 -hls_list_size 5 -hls_flags delete_segments -start_number 0 testsrc.m3u8
```

**Start HTTP server to serve HLS stream**

Reading the file m3u8 directly on the disk will prevent m3u8.tmp to be swapped with the main one. So the stream stops when reaching the last segments of the main one.

```sh
python3 -m http.server 8080
```

**Read the HLS stream**

```sh
.\ffplay -x 640 -y 360 http://127.0.0.1:8080/testsrc.m3u8
```

## Start a RTPM stream

Use `testsrc` source.

```sh
ffmpeg  -re -f lavfi -i testsrc -f flv rtmp://127.0.0.1:1935/live/app
```

Or use a playlist with `ffconcat`

```sh
ffmpeg  -re -stream_loop -1 -f concat -i playlist.txt -f flv rtmp://127.0.0.1:1935/live/app
```
