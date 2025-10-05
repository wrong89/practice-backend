package http

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
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
	router := mux.NewRouter()

	router.Path("/hello").Methods("GET").HandlerFunc(h.httpHandlers.HelloWorldHandler)
	router.Path("/user/register").Methods("POST").HandlerFunc(h.httpHandlers.RegisterUserHandler)

	return router
}
