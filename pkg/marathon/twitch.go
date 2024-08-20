package marathon

import (
	"errors"
	"strings"

	irc "github.com/gempir/go-twitch-irc/v4"
	log "github.com/sirupsen/logrus"
)

type TwitchConfig struct {
	TwitchOAuthToken              string
	TwitchStreamKey               string
	TwitchBroadcasterChannel      string
	TwitchWebsocketServer         string
	TwitchIrcServer               string
	TwitchAdminBadges             []string
	TwitchIngestServerUrlTemplate string
}

type TwitchService struct {
	Config *TwitchConfig

	streamingBroadcaster IStreamingBroadcaster
	ircClient            *irc.Client
}

func DefaultTwitchConfig() (s *TwitchConfig) {
	return &TwitchConfig{
		TwitchAdminBadges:             strings.Split(TwitchAdminBadges, ","),
		TwitchBroadcasterChannel:      "",
		TwitchIngestServerUrlTemplate: TwitchIngestServerUrlTemplate,
		TwitchIrcServer:               TwitchIrcServer,
		TwitchWebsocketServer:         TwitchWebsocketServer,
	}
}

func NewTwitchService(config *TwitchConfig) (IStreamingService, error) {

	if config.TwitchBroadcasterChannel == "" {
		return nil, errors.New("missing twitch broadcaster channel")
	}

	if config.TwitchStreamKey == "" {
		// Get Stream Key
		// channel:read:stream_key
		// https://dev.twitch.tv/docs/api/reference/#get-stream-key
		return nil, errors.New("missing twitch stream key")
	}

	if config.TwitchOAuthToken == "" {
		// Get OAuth Token
		// https://dev.twitch.tv/docs/cli/
		// https://dev.twitch.tv/docs/cli/token-command/
		log.Warn("missing twitch oauth token - bot will not interact in the channel, only listens to commands")
	}

	return &TwitchService{
		Config: config,
	}, nil
}

func (s *TwitchService) IsAdminBadge(badges map[string]int) bool {
	for badge := range badges {
		for _, adminBadge := range s.Config.TwitchAdminBadges {
			if badge == adminBadge {
				return true
			}
		}
	}
	return false
}

// Validate if Message is:
// - from same channel
// - from an allowed chatter
// - is prefixed with "!"
func (s *TwitchService) IsValidMessage(message irc.PrivateMessage) bool {
	if message.Channel != s.Config.TwitchBroadcasterChannel {
		// Wrong Channel
		log.Debugf("message ignored - wrong channel - %s", message.Channel)
		return false
	}
	if !s.IsAdminBadge(message.User.Badges) {
		// Not Admin
		log.Debugf("message ignored - no admin badges - %v", message.User.Badges)
		return false
	}
	// Validate Command
	if !strings.HasPrefix(message.Message, CommandPrefix) {
		// Not a command
		log.Debugf("message ignored - no command prefix - %s", message.Message)
		return false
	}
	return true
}

// Parse message and execute command
func (s *TwitchService) OnMessage(message string) {
	if PlayCommandRegexp.MatchString(message) {
		matches := PlayCommandRegexp.FindStringSubmatch(message)
		index := PlayCommandRegexp.SubexpIndex("path")
		if index < 0 || index > len(matches) {
			log.Warnf("message ignored - unexpected regexp match - %s with matches %v and index %d", message, matches, index)
		}
		path := matches[index]
		s.streamingBroadcaster.Play(path, s.StreamingUrl())
		return
	}
	if ShutdownCommandRegexp.MatchString(message) {
		s.Shutdown()
		return
	}
	log.Warnf("message ignored - unknown command - %s", message)
}

func (s *TwitchService) SetBroadcaster(streamingBroadcaster IStreamingBroadcaster) {
	s.streamingBroadcaster = streamingBroadcaster
}

// Depart from channel and disconnect from iRc server
func (s *TwitchService) Shutdown() {
	// TODO Replan the end of life cycle with multiple broadcasters and services
	s.streamingBroadcaster.Shutdown()
	if s.ircClient != nil {
		s.ircClient.Depart(s.Config.TwitchBroadcasterChannel)
		s.ircClient.Disconnect()
		s.ircClient = nil
	}
}

// Start
// - Join irs channel
// - Wait for commands
//   - Play -> kill/start stream video in loop
//   - Shutdown -> kill stream and bot
func (s *TwitchService) Start() {
	// https://dev.twitch.tv/docs/irc/#connecting-to-the-twitch-irc-server
	var ircClient *irc.Client
	if s.Config.TwitchOAuthToken != "" {
		ircClient = irc.NewClient(s.Config.TwitchBroadcasterChannel, s.Config.TwitchOAuthToken)
	} else {
		ircClient = irc.NewAnonymousClient()
	}
	s.ircClient = ircClient
	// https://dev.twitch.tv/docs/irc/capabilities/
	// TODO: Both are needed to have badges. Check to externalize access management.
	ircClient.Capabilities = []string{irc.TagsCapability, irc.CommandsCapability}
	ircClient.IrcAddress = s.Config.TwitchIrcServer
	ircClient.TLS = true

	// https://dev.twitch.tv/docs/irc/tags/#privmsg-tags
	ircClient.OnPrivateMessage(func(message irc.PrivateMessage) {
		if !s.IsValidMessage(message) {
			return
		}

		s.OnMessage(message.Message)
	})

	log.Infof("joining channel [%s]. waiting for commands...", s.Config.TwitchBroadcasterChannel)
	ircClient.Join(s.Config.TwitchBroadcasterChannel)

	err := ircClient.Connect()
	// Handles client disconnect
	if err == irc.ErrClientDisconnected {
		return
	}
	if err != nil {
		log.Errorf("failed to connect %v", err)
	}
}

// Build Streaming URL
func (s *TwitchService) StreamingUrl() string {
	return strings.Replace(s.Config.TwitchIngestServerUrlTemplate, "{stream_key}", s.Config.TwitchStreamKey, 1)
}
