package marathon

// https://dev.twitch.tv/docs/irc/#connecting-to-the-twitch-irc-server
const TwitchWebsocketServer string = "irc-ws.chat.twitch.tv:443"
const TwitchIrcServer string = "irc.chat.twitch.tv:6697"
const TwitchAdminBadges string = "admin,broadcaster,moderator"

// https://help.twitch.tv/s/twitch-ingest-recommendation
// https://dev.twitch.tv/docs/video-broadcast/reference/#get-ingest-servers
const TwitchIngestServerUrlTemplate string = "rtmps://cdg02.contribute.live-video.net/app/{stream_key}"

// https://help.twitch.tv/s/article/broadcasting-guidelines
// Example FFmpegArgs for 720p 30fps recommendation
// Change Video frame rate to 30, group of picture to 60 (2x frame rate), size 1280x760 (720p)
//   -framerate 30 -g 60 -video_size 1280x760
// Encode video with x264, preset veryfast (blind trust), video bit rate 3000kbps (max: 3000kbps), buffer size (< your bandwidth, > max bit rate, highly depend of your setup)
//   -c:v libx264 -preset veryfast -b:v 3000k -maxrate 3000k -bufsize 1M
// Encode audio with AAC, audio bit rate 128kbps (reco: 96, max: 160), audio sampling rate 44100Hz
//   -c:a aac -b:a 128k -ar 44100
const FFmpegArgs string = "-hide_banner -loglevel error -framerate 30 -g 60 -video_size 1280x760 -c:v libx264 -preset veryfast -b:v 3000k -maxrate 3000k -bufsize 1M -c:a aac -b:a 128k -ar 44100"

const StreamerTypeFFmpeg = "ffmpeg"
