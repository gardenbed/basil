package main

import (
	"context"
	"errors"
	"time"
)

// Client is a mock client implementing graceful.Client interface.
type Client struct {
	count int
	name  string
}

// NewClient creates a new client.
func NewClient(name string) *Client {
	return &Client{
		name: name,
	}
}

func (c *Client) String() string {
	return c.name
}

// Connect mimicks connecting to an external service.
func (c *Client) Connect() error {
	c.count++
	time.Sleep(time.Second)

	// For testing retries
	if c.count < 3 {
		return errors.New("error on connecting client")
	}

	return nil
}

// Disconnect mimicks disconnecting from an external service.
func (c *Client) Disconnect(context.Context) error {
	time.Sleep(time.Second)
	return nil
}
