package ghttp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net"
	"net/http"
	"strconv"
)

type Server struct {
	Address string
	Port    int
	engine  *gin.Engine
}

func NewServer() *Server {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery(), GinLogger())
	binding.Validator = new(Validator)

	return &Server{engine: engine}
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
