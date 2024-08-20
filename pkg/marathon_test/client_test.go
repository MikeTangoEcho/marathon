package marathon_test

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/MikeTangoEcho/marathon/pkg/marathon"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestNewClient(t *testing.T) {
	assert := assert.New(t)

	ss := &StreamingServiceMock{}
	sb := &StreamingBroadcasterMock{}

	ss.On("SetBroadcaster", mock.Anything).Return()

	c := marathon.NewClient(ss, sb)

	assert.NotNil(c)

	cc, ok := c.(*marathon.Client)
	assert.True(ok)

	assert.Same(cc.StreamingService, ss)
	assert.Same(cc.StreamingBroadcaster, sb)

	ss.AssertCalled(t, "SetBroadcaster", sb)
}

func TestClient_Run(t *testing.T) {
	assert := assert.New(t)

	ss := &StreamingServiceMock{}
	sb := &StreamingBroadcasterMock{}

	ss.On("Shutdown").Return()
	sb.On("Shutdown").Return()

	ss.On("SetBroadcaster", mock.Anything).Return()
	ss.On("Start").Return()

	c := marathon.NewClient(ss, sb)
	assert.NotNil(c)

	var call *mock.Call = nil

	t.Run("broadcaster failed to prepare", func(t *testing.T) {
		call = sb.On("Prepare").Return(errors.New("failed"))

		c.Run()

		sb.AssertCalled(t, "Prepare")
		ss.AssertNotCalled(t, "Start")

		ss.AssertCalled(t, "Shutdown")
		sb.AssertCalled(t, "Shutdown")
	})

	t.Run("success", func(t *testing.T) {
		call.Unset()
		call = sb.On("Prepare").Return(nil)

		c.Run()

		sb.AssertCalled(t, "Prepare")
		ss.AssertCalled(t, "Start")

		ss.AssertCalled(t, "Shutdown")
		sb.AssertCalled(t, "Shutdown")
	})
}
