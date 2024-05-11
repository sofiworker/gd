package gdrpc

import "google.golang.org/grpc"


func NewGRPCClient() (*grpc.ClientConn, error) {
	server, err := grpc.NewClient("")
	if err != nil {
		return nil, err
	}
	return server, nil
}
