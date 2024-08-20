package marathon_test

import (
	"io/ioutil"
	"testing"

	"github.com/MikeTangoEcho/marathon/pkg/marathon"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestDefaultFFmpegConfig(t *testing.T) {
	assert := assert.New(t)

	c := marathon.DefaultFFmpegConfig()
	assert.NotNil(c)
	assert.Empty(c.FFmpegPath)
}

func TestNewFFmpegBroadcaster(t *testing.T) {
	assert := assert.New(t)

	c := marathon.DefaultFFmpegConfig()
	b, err := marathon.NewFFmpegBroadcaster(c)
	assert.Nil(err)
	assert.NotNil(b)
}

func TestFFmpegBroadcaster_CommandArgs(t *testing.T) {
	assert := assert.New(t)

	c := marathon.DefaultFFmpegConfig()
	b, err := marathon.NewFFmpegBroadcaster(c)
	assert.Nil(err)
	assert.NotNil(b)

	sb, ok := b.(*marathon.FFmpegBroadcaster)
	assert.True(ok)

	// TODO mock filesystem/cmd to assert that the proper command has been built
	sb.Config.FFmpegArgs = "-other-args"
	args := sb.CommandArgs("file", "rtmp://")
	assert.Equal([]string{
		"-re",
		"-stream_loop", "1",
		"-f", "concat",
		"-i", "file",
		"-other-args",
		"-f", "flv",
		"rtmp://",
	}, args)
}

func TestFFmpegBroadcaster_Play(t *testing.T) {
	assert := assert.New(t)

	c := marathon.DefaultFFmpegConfig()
	b, err := marathon.NewFFmpegBroadcaster(c)
	assert.Nil(err)
	assert.NotNil(b)

	// TODO mock filesystem/cmd to assert that the proper command has been built
	b.Play("file", "rtmp://")
}

func TestFFmpegBroadcaster_Prepare(t *testing.T) {
	assert := assert.New(t)

	c := marathon.DefaultFFmpegConfig()
	b, err := marathon.NewFFmpegBroadcaster(c)
	assert.Nil(err)
	assert.NotNil(b)

	// TODO mock filesystem/cmd to assert that we lookup the proper version of ffmpeg
	err = b.Prepare()

	assert.NotNil(err)
}

func TestFFmpegBroadcaster_Shutdown(t *testing.T) {
	assert := assert.New(t)

	c := marathon.DefaultFFmpegConfig()
	b, err := marathon.NewFFmpegBroadcaster(c)
	assert.Nil(err)
	assert.NotNil(b)

	b.Shutdown()
}
