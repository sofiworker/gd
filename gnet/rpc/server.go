package rpc

import (
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	s *grpc.Server
}

func NewServer() *Server {
	server := grpc.NewServer()
	return &Server{s: server}
}

func (s *Server) Register(sd *grpc.ServiceDesc, ss any) {
	s.s.RegisterService(sd, ss)
}

func (s *Server) Start() error {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}
	err = s.s.Serve(listen)
	if err != nil {
		return err
	}
	return nil
}
