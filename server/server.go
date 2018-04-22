package server

import (
	"net/http"

	"github.com/go-chi/chi/middleware"

	"github.com/go-chi/chi"
)

// Server configures the HTTP server
type Server struct {
	Host string
}

// Run the HTTP server
func (s *Server) Run() error {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	sock := newSocketServer()
	r.HandleFunc("/api/socket", sock.socketUpgrader)
	r.Handle("/*", http.FileServer(http.Dir("website")))

	return http.ListenAndServe(s.Host, r)
}
