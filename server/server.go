package server

import "net/http"

type Server struct {
	httpServer *http.Server
}

func New(address string) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    address,
			Handler: http.NewServeMux(),
		},
	}
}

func (server *Server) ListenAndServe() error {
	return server.httpServer.ListenAndServe()
}
