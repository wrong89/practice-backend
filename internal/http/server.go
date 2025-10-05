package http

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HTTPServer struct {
	httpHandlers HTTPHandlers
}

func NewHTTPServer(httpHandlers HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		httpHandlers: httpHandlers,
	}
}

func (h *HTTPServer) Start() error {
	router := h.configureRouter()

	if err := http.ListenAndServe(":9091", router); err != nil {
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

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World PUBLIC"))
		})

		r.Post("/user/register", h.httpHandlers.RegisterHandler)
		r.Post("/user/login", h.httpHandlers.LoginHandler)

		r.With(AuthMiddleware).Get("/test", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello FROM PRIVATE ROUTE"))
		})
	})

	// router.Use(LoggingMiddleware)
	// router.Use(AuthMiddleware)

	// router.Path("/user/register").Methods("POST").HandlerFunc(h.httpHandlers.RegisterHandler)
	// router.Path("/user/login").Methods("POST").HandlerFunc(h.httpHandlers.LoginHandler)

	return router
}
