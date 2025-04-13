package rpc

import (
	"google.golang.org/grpc"
	"net"
)

type ServerOpts struct {
	DisallowUnknownFields bool
	OtelName              string
}

type ServerOptsFunc func(opts *ServerOpts)

func WithDisallowUnknownFields() ServerOptsFunc {
	return func(opts *ServerOpts) {
		opts.DisallowUnknownFields = true
	}
}

func WithTraceName(name string) ServerOptsFunc {
	return func(opts *ServerOpts) {
		opts.OtelName = name
	}
}

type Server struct {
	s    *grpc.Server
	opts *ServerOpts
}

func NewServer(opts ...ServerOptsFunc) *Server {
	server := grpc.NewServer()
	var sOpts ServerOpts
	for _, opt := range opts {
		opt(&sOpts)
	}
	return &Server{s: server, opts: &sOpts}
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
