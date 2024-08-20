package marathon

import (
	log "github.com/sirupsen/logrus"
)

func NewClient(streamingService IStreamingService, streamingBroadcaster IStreamingBroadcaster) IClient {
	streamingService.SetBroadcaster(streamingBroadcaster)
	return &Client{
		StreamingService:     streamingService,
		StreamingBroadcaster: streamingBroadcaster,
	}
}

func (c *Client) Halt() {
	c.StreamingService.Shutdown()
	c.StreamingBroadcaster.Shutdown()
}

func (c *Client) Run() {
	defer log.Info("bye 👋")
	defer c.Halt()
	log.Info("starting marathon 🏃...")

	err := c.StreamingBroadcaster.Prepare()
	if err != nil {
		log.Error(err)
		return
	}
	c.StreamingService.Start()
}
