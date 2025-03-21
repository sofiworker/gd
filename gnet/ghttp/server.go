package ghttp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"net"
	"net/http"
	"strconv"
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
	Address string
	Port    int
	opts    *ServerOpts
	engine  *gin.Engine
}

func NewServer(opts ...ServerOptsFunc) *Server {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(gin.Recovery(), GinLogger())

	binding.Validator = new(Validator)

	var sOpts ServerOpts
	for _, opt := range opts {
		opt(&sOpts)
	}

	if sOpts.DisallowUnknownFields {
		binding.EnableDecoderDisallowUnknownFields = true
	}

	if sOpts.OtelName == "" {
		sOpts.OtelName = "gin"
	}
	engine.Use(otelgin.Middleware(sOpts.OtelName))

	return &Server{engine: engine, opts: &sOpts}
}

// ListenAndServe you should use this method in goroutine and then use gd.ListenShutDownSignals
func (s *Server) ListenAndServe() error {
	host := net.JoinHostPort(s.Address, strconv.Itoa(s.Port))
	if err := http.ListenAndServe(host, s.engine); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
	return nil
}

func (s *Server) Post(path string, handler ...gin.HandlerFunc) {
	s.engine.POST(path, handler...)
}

func (s *Server) Head(path string, handler ...gin.HandlerFunc) {
	s.engine.HEAD(path, handler...)
}

func (s *Server) Get(path string, handler ...gin.HandlerFunc) {
	s.engine.GET(path, handler...)
}

func (s *Server) Put(path string, handler ...gin.HandlerFunc) {
	s.engine.PUT(path, handler...)
}

func (s *Server) Delete(path string, handler ...gin.HandlerFunc) {
	s.engine.DELETE(path, handler...)
}

func (s *Server) Patch(path string, handler ...gin.HandlerFunc) {
	s.engine.PATCH(path, handler...)
}

func (s *Server) Any(path string, handler ...gin.HandlerFunc) {
	s.engine.Any(path, handler...)
}

func (s *Server) StaticFile(relativePath, filepath string) {
	s.engine.StaticFile(relativePath, filepath)
}

func (s *Server) Static(relativePath, root string) {
	s.engine.Static(relativePath, root)
}
