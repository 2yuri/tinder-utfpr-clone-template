package api

import (
	"context"
	"github.com/gin-gonic/gin"
)

type Server struct {
	context     context.Context
	port        string
	app         *gin.Engine
	releaseMode bool
}

type ServerOptions func(*Server) error

func NewServer(ctx context.Context, port string, debugMode bool) *Server {
	server := &Server{context: ctx, port: port, app: gin.New(), releaseMode: debugMode}

	return server
}
func (s *Server) Port() string {
	return s.port
}

func (s *Server) Start() *gin.Engine {
	if s.releaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	s.setupDefaultMiddlewares()
	s.setupRoutes()

	return s.app
}

func (s *Server) setupDefaultMiddlewares() {
	s.app.Use(gin.Logger())
	s.app.Use(gin.Recovery())
}
