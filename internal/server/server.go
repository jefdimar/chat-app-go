package server

import (
	"net/http"

	"github.com/jefdimar/go-chat-app/internal/handlers"
)

// Server represents the chat server
type Server struct {
	Router *http.ServeMux
}

// NewServer creates a new chat server
func NewServer() *Server {
	s := &Server{
		Router: http.NewServeMux(),
	}

	// Register routes
	s.routes()

	return s
}

// routes sets up the routes for the server
func (s *Server) routes() {
	// Serve static files
	fileServer := http.FileServer(http.Dir("./web/static"))
	s.Router.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Home Page
	s.Router.HandleFunc("/", handlers.Home)

	// WebSocket endpoint
	s.Router.HandleFunc("/ws", handlers.HandleConnections)
}
