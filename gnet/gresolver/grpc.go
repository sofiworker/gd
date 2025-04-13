package gresolver

import (
	"context"
	"google.golang.org/grpc/resolver"
	"sync"
)

type dnsBuilder struct{}

// NewBuilder creates a dnsBuilder which is used to factory DNS resolvers.
func NewBuilder() resolver.Builder {
	return &dnsBuilder{}
}

// Build creates and starts a DNS resolver that watches the name resolution of
// the target.
func (b *dnsBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	return nil, nil
}

// Scheme returns the naming scheme of this resolver builder, which is "dns".
func (b *dnsBuilder) Scheme() string {
	return "dns"
}

// dnsResolver watches for the name resolution update for a non-IP target.
type dnsResolver struct {
	host string
	port string
	//resolver internal.NetResolver
	ctx    context.Context
	cancel context.CancelFunc
	cc     resolver.ClientConn
	// rn channel is used by ResolveNow() to force an immediate resolution of the
	// target.
	rn chan struct{}
	// wg is used to enforce Close() to return after the watcher() goroutine has
	// finished. Otherwise, data race will be possible. [Race Example] in
	// dns_resolver_test we replace the real lookup functions with mocked ones to
	// facilitate testing. If Close() doesn't wait for watcher() goroutine
	// finishes, race detector sometimes will warns lookup (READ the lookup
	// function pointers) inside watcher() goroutine has data race with
	// replaceNetFunc (WRITE the lookup function pointers).
	wg                   sync.WaitGroup
	disableServiceConfig bool
}

// ResolveNow invoke an immediate resolution of the target that this
// dnsResolver watches.
func (d *dnsResolver) ResolveNow(resolver.ResolveNowOptions) {
	select {
	case d.rn <- struct{}{}:
	default:
	}
}

// Close closes the dnsResolver.
func (d *dnsResolver) Close() {
}

func RegisterGrpcResolver(b resolver.Builder) error {
	builder := resolver.Get(b.Scheme())
	if builder != nil {
		return nil
	}
	resolver.Register(b)
	return nil
}
