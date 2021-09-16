package main

import (
	"context"
	"time"
)

// Client is a mock client implementing health.Checker interface.
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

// HealthCheck mimicks the health checking logic of the client.
func (c *Client) HealthCheck(context.Context) error {
	time.Sleep(100 * time.Millisecond)
	return nil
}
