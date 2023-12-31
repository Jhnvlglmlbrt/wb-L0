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

func (s *Server) Start(handlePage echo.HandlerFunc, handleGetData echo.HandlerFunc) error {
	s.echo.GET("/", handlePage)
	s.echo.GET("/order", handleGetData)
	s.echo.Static("/", "static")

	return s.echo.Start(s.port)
}
