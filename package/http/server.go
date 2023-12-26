package http

import (
	"github.com/Jhnvlglmlbrt/wb-order/config"
	"github.com/labstack/echo/v4"
)

type Server struct {
	echo *echo.Echo
	port string
}

func NewServer(cfg *config.HTTP) *Server {
	return &Server{
		echo: echo.New(),
		port: ":" + cfg.Port,
	}
}

func (s *Server) Start(orderHandler echo.HandlerFunc, allOrdersHandler echo.HandlerFunc) error {
	s.echo.GET("/:order", orderHandler)
	s.echo.GET("/", allOrdersHandler)
	return s.echo.Start(s.port)
}
