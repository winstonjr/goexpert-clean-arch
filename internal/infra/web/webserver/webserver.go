package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router        chi.Router
	Handlers      []BindInformation
	WebServerPort string
}

type BindInformation struct {
	Path         string
	Method       string
	RouteHandler http.HandlerFunc
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      []BindInformation{},
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(verb string, path string, handler http.HandlerFunc) {
	s.Handlers = append(s.Handlers, BindInformation{path, verb, handler})
}

// loop through the handlers and add them to the router
// register middeleware logger
// start the server
func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	for _, handler := range s.Handlers {
		if handler.Method == http.MethodGet {
			s.Router.Get(handler.Path, handler.RouteHandler)
		} else if handler.Method == http.MethodPost {
			s.Router.Post(handler.Path, handler.RouteHandler)
		} else {
			s.Router.Handle(handler.Path, handler.RouteHandler)
		}
	}
	http.ListenAndServe(s.WebServerPort, s.Router)
}
