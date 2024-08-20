package marathon_test

import (
	"io/ioutil"
	"testing"

	"github.com/MikeTangoEcho/marathon/pkg/marathon"
	irc "github.com/gempir/go-twitch-irc/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func DummyTwitchConfig() *marathon.TwitchConfig {
	c := marathon.DefaultTwitchConfig()
	c.TwitchBroadcasterChannel = "justinfanxxx"
	c.TwitchStreamKey = "live_xxx"
	return c
}

func TestDefaultTwitchConfig(t *testing.T) {
	assert := assert.New(t)

	c := marathon.DefaultTwitchConfig()
	assert.NotNil(c)
	assert.Empty(c.TwitchBroadcasterChannel)
}

func TestNewTwitchService(t *testing.T) {
	assert := assert.New(t)

	c := marathon.DefaultTwitchConfig()

	t.Run("missing TwitchBroadcasterChannel", func(t *testing.T) {
		c.TwitchBroadcasterChannel = ""
		s, err := marathon.NewTwitchService(c)
		assert.Error(err)
		assert.Nil(s)
	})
	t.Run("missing TwitchStreamKey", func(t *testing.T) {
		c.TwitchBroadcasterChannel = "justinfanxxx"
		c.TwitchStreamKey = ""
		s, err := marathon.NewTwitchService(c)
		assert.Error(err)
		assert.Nil(s)
	})
	t.Run("success", func(t *testing.T) {
		c.TwitchBroadcasterChannel = "justinfanxxx"
		c.TwitchStreamKey = "live_xxx"
		s, err := marathon.NewTwitchService(c)
		assert.Nil(err)
		assert.NotNil(s)
	})
}

func TestTwitchService_IsAdminBadge(t *testing.T) {
	assert := assert.New(t)

	s, err := marathon.NewTwitchService(DummyTwitchConfig())
	assert.Nil(err)
	assert.NotNil(s)

	ss, ok := s.(*marathon.TwitchService)
	assert.True(ok)

	assert.True(ss.IsAdminBadge(map[string]int{
		"admin": 0,
	}))
	assert.True(ss.IsAdminBadge(map[string]int{
		"broadcaster": 0,
	}))
	assert.True(ss.IsAdminBadge(map[string]int{
		"moderator": 0,
	}))
	assert.False(ss.IsAdminBadge(map[string]int{
		"subscriber": 42,
	}))
	ss.Config.TwitchAdminBadges = append(ss.Config.TwitchAdminBadges, "subscriber")
	assert.True(ss.IsAdminBadge(map[string]int{
		"subscriber": 42,
	}))
}

func TestTwitchService_IsValidMessage(t *testing.T) {
	assert := assert.New(t)

	c := DummyTwitchConfig()
	s, err := marathon.NewTwitchService(c)
	assert.Nil(err)
	assert.NotNil(s)

	ss, ok := s.(*marathon.TwitchService)
	assert.True(ok)
	t.Run("invalid broadcast channel", func(t *testing.T) {
		assert.False(ss.IsValidMessage(irc.PrivateMessage{
			Channel: c.TwitchBroadcasterChannel + "wtv",
		}))
	})

	t.Run("invalid badges", func(t *testing.T) {
		assert.False(ss.IsValidMessage(irc.PrivateMessage{
			Channel: c.TwitchBroadcasterChannel,
			User: irc.User{
				Badges: map[string]int{
					"subscriber": 42,
				},
			},
		}))
	})

	assert.NotEmpty(c.TwitchAdminBadges)
	var adminBadges map[string]int = make(map[string]int)
	for _, badge := range c.TwitchAdminBadges {
		adminBadges[badge] = 0
	}

	t.Run("invalid command", func(t *testing.T) {
		assert.False(ss.IsValidMessage(irc.PrivateMessage{
			Channel: c.TwitchBroadcasterChannel,
			User: irc.User{
				Badges: adminBadges,
			},
			Message: "notacommand",
		}))
	})

	t.Run("success", func(t *testing.T) {
		assert.True(ss.IsValidMessage(irc.PrivateMessage{
			Channel: c.TwitchBroadcasterChannel,
			User: irc.User{
				Badges: adminBadges,
			},
			Message: "!command",
		}))
	})
}

func TestTwitchService_OnMessage(t *testing.T) {
	assert := assert.New(t)

	s, err := marathon.NewTwitchService(DummyTwitchConfig())
	assert.Nil(err)
	assert.NotNil(s)

	sbsm := &StreamingBroadcasterMock{}
	sbsm.On("Play", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return()
	sbsm.On("Shutdown").Return()

	s.SetBroadcaster(sbsm)

	t.Run("not a command", func(t *testing.T) {
		s.OnMessage("notacommand")
		sbsm.AssertNotCalled(t, "Shutdown")
		sbsm.AssertNotCalled(t, "Play")
		// TODO add helper to asset all available command
	})

	t.Run("unkown command", func(t *testing.T) {
		s.OnMessage("!unkowncommand")
		sbsm.AssertNotCalled(t, "Play")
		sbsm.AssertNotCalled(t, "Shutdown")
		// TODO add helper to asset all available command
	})

	t.Run("malformed command", func(t *testing.T) {
		s.OnMessage("!playtest")
		sbsm.AssertNotCalled(t, "Play")
		sbsm.AssertNotCalled(t, "Shutdown")
		s.OnMessage("!play something and other uneeded args")
		sbsm.AssertNotCalled(t, "Play")
		sbsm.AssertNotCalled(t, "Shutdown")
		s.OnMessage("!shutdown now")
		sbsm.AssertNotCalled(t, "Play")
		sbsm.AssertNotCalled(t, "Shutdown")
	})

	t.Run("success play", func(t *testing.T) {
		s.OnMessage("!play somewhere/something\\elsewhere")
		sbsm.AssertCalled(t, "Play", "somewhere/something\\elsewhere", s.StreamingUrl())
		sbsm.AssertNotCalled(t, "Shutdown")
	})

	t.Run("success shutdown", func(t *testing.T) {
		s.OnMessage("!shutdown")
		sbsm.AssertNotCalled(t, "Play")
		sbsm.AssertCalled(t, "Shutdown")
	})
}

func TestTwitchService_SetBroadcaster(t *testing.T) {
	assert := assert.New(t)

	s, err := marathon.NewTwitchService(DummyTwitchConfig())
	assert.Nil(err)
	assert.NotNil(s)

	sbsm := &StreamingBroadcasterMock{}
	s.SetBroadcaster(sbsm)
	_, ok := s.(*marathon.TwitchService)
	assert.True(ok)

	// TODO private
}

func TestTwitchService_Shutdown(t *testing.T) {
	assert := assert.New(t)

	s, err := marathon.NewTwitchService(DummyTwitchConfig())
	assert.Nil(err)
	assert.NotNil(s)

	ss, ok := s.(*marathon.TwitchService)
	assert.True(ok)

	sbsm := &StreamingBroadcasterMock{}
	sbsm.On("Shutdown").Return()

	ss.SetBroadcaster(sbsm)
	ss.Shutdown()

	sbsm.AssertCalled(t, "Shutdown")
}

func TestTwitchService_Start(t *testing.T) {
	assert := assert.New(t)

	s, err := marathon.NewTwitchService(DummyTwitchConfig())
	assert.Nil(err)
	assert.NotNil(s)

	// TODO mock iRc server
	//s.Start()
}

func TestTwitchService_StreamingUrl(t *testing.T) {
	assert := assert.New(t)

	s, err := marathon.NewTwitchService(DummyTwitchConfig())
	assert.Nil(err)
	assert.NotNil(s)

	url := s.StreamingUrl()

	assert.NotEmpty(url)
	assert.NotContains(url, "{stream_key}")
}
