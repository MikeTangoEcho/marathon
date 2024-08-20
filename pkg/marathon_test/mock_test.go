package marathon_test

import (
	"github.com/MikeTangoEcho/marathon/pkg/marathon"
	"github.com/stretchr/testify/mock"
)

type StreamingBroadcasterMock struct {
	mock.Mock
}

func (s *StreamingBroadcasterMock) Prepare() error {
	args := s.Called()
	return args.Error(0)
}

func (s *StreamingBroadcasterMock) Shutdown() {
	s.Called()
}

func (s *StreamingBroadcasterMock) Play(args string, streamingUrl string) {
	s.Called(args, streamingUrl)
}

type StreamingServiceMock struct {
	mock.Mock
}

func (s *StreamingServiceMock) OnMessage(message string) {
	s.Called(message)
}

func (s *StreamingServiceMock) SetBroadcaster(broadcaster marathon.IStreamingBroadcaster) {
	s.Called(broadcaster)
}

func (s *StreamingServiceMock) Shutdown() {
	s.Called()
}

func (s *StreamingServiceMock) Start() {
	s.Called()
}

func (s *StreamingServiceMock) StreamingUrl() string {
	args := s.Called()
	return args.String()
}
