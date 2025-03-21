package ghttp

import (
	"context"
	"net"
)

type Resolver interface {
	//Resolve(string) ([]net.Addr, error)
	//ResolveContext(ctx context.Context, name string) ([]net.Addr, error)
	GoResolve(ctx context.Context, network, address string) (net.Conn, error)
}

type DefaultResolver struct {
	remote []string
}

func NewDefaultResolver(remote ...string) *DefaultResolver {
	return &DefaultResolver{remote: remote}
}

//func (r *DefaultResolver) Resolve(name string) ([]net.Addr, error) {
//	return r.ResolveContext(context.Background(), name)
//}
//
//func (r *DefaultResolver) ResolveContext(ctx context.Context, name string) ([]net.Addr, error) {
//	if name == "" {
//		return nil, NotFoundMethodError
//	}
//	return nil, nil
//}

func (r *DefaultResolver) GoResolve(ctx context.Context, network, address string) (net.Conn, error) {
	return nil, nil
}
