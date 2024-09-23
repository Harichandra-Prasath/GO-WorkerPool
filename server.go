package main

import (
	"fmt"
	"net/http"
)

type ServerConfig struct {
	Addr string
}

type Server struct {
	Config *ServerConfig
	Pool   *Pool
	Mux    *http.ServeMux
}

func GetNewServer(cfg *ServerConfig) *Server {
	return &Server{
		Config: cfg,
		Mux:    http.NewServeMux(),
	}
}

func (S *Server) Serve() error {

	fmt.Printf("Server Started on localhost:%s\n", S.Config.Addr)
	err := http.ListenAndServe(S.Config.Addr, S.Mux)
	if err != nil {
		return err
	}
	return nil
}
