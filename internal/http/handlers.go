package http

import (
	"context"
	"encoding/json"
	"net/http"
	"practice-backend/internal/models/entry"
	"practice-backend/internal/models/user"
	"practice-backend/internal/services/auth"
)

type Auth interface {
	Login(
		ctx context.Context,
		login string,
		password string,
	) (token string, err error)
	Register(
		ctx context.Context,
	) (userID int64, err error)
}

type HTTPHandlers struct {
	entryRepo   entry.EntryRepo
	userRepo    user.UserRepo
	authService auth.Auth
}

func NewHTTPHandlers(
	entryRepo entry.EntryRepo,
	userRepo user.UserRepo,
	authService auth.Auth,
) *HTTPHandlers {
	return &HTTPHandlers{
		entryRepo:   entryRepo,
		userRepo:    userRepo,
		authService: authService,
	}
}

func (h *HTTPHandlers) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	var helloWorldDTO HelloWorldDTO

	if err := json.NewDecoder(r.Body).Decode(&helloWorldDTO); err != nil {
		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.String(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("HelloWorld for:\t" + helloWorldDTO.Title))
}

/* contract example
pattern: /user/register
method:  POST
info:    JSON in HTTP request body

succeed:
  - status code: 201 Created
  - response body: JSON of created created user
failed:
  - status code: 400, 409(Conflict), 500
  - response body: JSON with error + time
*/

func (h *HTTPHandlers) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var helloWorldDTO HelloWorldDTO

	if err := json.NewDecoder(r.Body).Decode(&helloWorldDTO); err != nil {
		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.String(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("HelloWorld for:\t" + helloWorldDTO.Title))
}
