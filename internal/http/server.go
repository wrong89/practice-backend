package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HTTPServer struct {
	httpHandlers HTTPHandlers
	port         int
	host         string
}

func NewHTTPServer(httpHandlers HTTPHandlers, port int, host string) *HTTPServer {
	return &HTTPServer{
		httpHandlers: httpHandlers,
		port:         port,
		host:         host,
	}
}

func (h *HTTPServer) Start() error {
	router := h.configureRouter()

	socket := fmt.Sprintf("%s:%d", h.host, h.port)

	if err := http.ListenAndServe(socket, router); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	return nil
}

func (h *HTTPServer) configureRouter() http.Handler {
	router := chi.NewRouter()

	router.Route("/api", func(r chi.Router) {
		r.Use(LoggingMiddleware)
		r.Use(CorsMiddleware)

		r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("HelloWorld"))
		})

		r.Post("/user/register", h.httpHandlers.RegisterHandler)
		r.Post("/user/login", h.httpHandlers.LoginHandler)

		r.With(AuthMiddleware).Post("/entry", h.httpHandlers.CreateEntryHandler)
		r.With(AuthMiddleware).Get("/entry", h.httpHandlers.GetEntriesHandler)
		r.With(AuthMiddleware).Patch("/entry", h.httpHandlers.UpdateEntryHandler)

		r.With(AdminMiddleware(h.httpHandlers.authService)).Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ADMIN WRITE"))
		})
	})

	return router
}
