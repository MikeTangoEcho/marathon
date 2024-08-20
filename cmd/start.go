package cmd

import (
	"log"

	"github.com/MikeTangoEcho/marathon/pkg/marathon"
	"github.com/spf13/cobra"
)

const TwitchCmdFlagValue = "twitch"
const FFmpegCmdFlagValue = "ffmpeg"

var StreamingService string
var StreamingBroadcaster string
var TwitchConfig *marathon.TwitchConfig = marathon.DefaultTwitchConfig()
var FFmpegConfig *marathon.FFmpegConfig = marathon.DefaultFFmpegConfig()

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the bot",
	Long:  ``,
	Args:  cobra.NoArgs,
	PreRun: func(cmd *cobra.Command, args []string) {
		if StreamingService == TwitchCmdFlagValue {
			cmd.MarkFlagRequired("twitch-broadcaster-channel")
			cmd.MarkFlagRequired("twitch-stream-key")
		}
		if StreamingBroadcaster == FFmpegCmdFlagValue {

		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var s marathon.IStreamingService
		var b marathon.IStreamingBroadcaster
		var err error

		switch StreamingService {
		case TwitchCmdFlagValue:
			s, err = marathon.NewTwitchService(TwitchConfig)
			if err != nil {
				panic(err)
			}
		default:
			log.Panicf("unkown streaming service %s", StreamingService)

		}

		switch StreamingBroadcaster {
		case FFmpegCmdFlagValue:
			b, err = marathon.NewFFmpegBroadcaster(FFmpegConfig)
			if err != nil {
				panic(err)
			}
		default:
			log.Panicf("unkown streaming broadcaster %s", StreamingBroadcaster)
		}

		c := marathon.NewClient(s, b)
		c.Run()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Type
	startCmd.Flags().StringVar(&StreamingService, "streaming-service", TwitchCmdFlagValue, "Streaming service type.")
	startCmd.Flags().StringVar(&StreamingBroadcaster, "streaming-broadcaster", FFmpegCmdFlagValue, "Streaming broacaster type.")

	// Twitch
	startCmd.Flags().StringVar(&TwitchConfig.TwitchOAuthToken, "twitch-oauth-token", "", "Twitch OAuth Token used for Chat Bot.")
	startCmd.Flags().StringVar(&TwitchConfig.TwitchStreamKey, "twitch-stream-key", "", "Twitch Stream Key used to stream to Twitch.")
	startCmd.Flags().StringArrayVar(&TwitchConfig.TwitchAdminBadges, "twitch-admin-badges", TwitchConfig.TwitchAdminBadges, "")
	startCmd.Flags().StringVar(&TwitchConfig.TwitchBroadcasterChannel, "twitch-broadcaster-channel", "", "Twitch Broadcaster Channel that the Bot join to listen to commands.")
	startCmd.Flags().StringVar(&TwitchConfig.TwitchIngestServerUrlTemplate, "twitch-ingest-server-url-template", TwitchConfig.TwitchIngestServerUrlTemplate, "Twitch Ingest Server Url Template. https://help.twitch.tv/s/twitch-ingest-recommendation")
	startCmd.Flags().StringVar(&TwitchConfig.TwitchIrcServer, "twitch-irc-server", TwitchConfig.TwitchIrcServer, "Twitch iRc Server")
	startCmd.Flags().StringVar(&TwitchConfig.TwitchWebsocketServer, "twitch-websocket-server", TwitchConfig.TwitchWebsocketServer, "Twitch Websocket Server")

	// FFmpeg
	startCmd.Flags().StringVar(&FFmpegConfig.FFmpegArgs, "ffmpeg-args", FFmpegConfig.FFmpegArgs, "FFmpeg Args used to encode. Can be omitted if your videos don't need reencoding. https://help.twitch.tv/s/article/broadcasting-guidelines")
	startCmd.Flags().StringVar(&FFmpegConfig.FFmpegPath, "ffmpeg-path", "", "FFmpeg bin path. Can be relative. If omitted, marathon will use the one in PATH, else download it https://www.ffmpeg.org/download.html")
}
