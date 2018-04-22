package server

import (
	"net/http"

	"github.com/go-chi/chi/middleware"

	"github.com/go-chi/chi"
)

// Server configures the HTTP server
type Server struct {
	Host   string
	socket *socketServer
}

// Run the HTTP server
func (s *Server) Run() error {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	sock := newSocketServer()
	go sock.run()
	s.socket = sock

	r.HandleFunc("/api/socket", sock.socketUpgrader)
	r.HandleFunc("/api/webhook/{key}", s.StatusWebhook)
	r.Handle("/*", http.FileServer(http.Dir("website")))

	return http.ListenAndServe(s.Host, r)
}

// StatusWebhook handles incoming status webhooks from the bots
func (s *Server) StatusWebhook(w http.ResponseWriter, r *http.Request) {

}
