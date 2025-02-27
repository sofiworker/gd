package rpc

import "google.golang.org/grpc"

type Client struct {
	c *grpc.ClientConn
}

func NewClient() (*Client, error) {
	c, err := grpc.NewClient("")
	if err != nil {
		return nil, err
	}
	return &Client{
		c: c,
	}, nil
}
