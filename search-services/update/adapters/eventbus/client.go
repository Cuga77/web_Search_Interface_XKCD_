package eventbus

import (
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"
)

type Client struct {
	nc  *nats.Conn
	log *slog.Logger
}

func NewClient(brokerAddress string, log *slog.Logger) (*Client, error) {
	if brokerAddress == "" {
		return nil, fmt.Errorf("broker address is empty")
	}
	nc, err := nats.Connect(brokerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to nats: %w", err)
	}
	return &Client{nc: nc, log: log}, nil
}

func (c *Client) PublishUpdate() error {
	c.log.Debug("publishing xkcd.db.updated event")
	err := c.nc.Publish("xkcd.db.updated", []byte("updated"))
	if err != nil {
		return fmt.Errorf("failed to publish update event: %w", err)
	}
	return c.nc.Flush()
}

func (c *Client) Close() {
	if c.nc != nil {
		c.nc.Close()
	}
}
