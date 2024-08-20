package marathon

type IClient interface {
	Run()
}

type Client struct {
	StreamingBroadcaster IStreamingBroadcaster
	StreamingService     IStreamingService
}

type IStreamingBroadcaster interface {
	Prepare() error
	Play(path string, streamingUrl string)
	Shutdown()
}

type IStreamingService interface {
	SetBroadcaster(broadcaster IStreamingBroadcaster)
	OnMessage(message string)
	Shutdown()
	Start()
	StreamingUrl() string
}
