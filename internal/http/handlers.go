package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"practice-backend/internal/models/entry"
	"practice-backend/internal/models/user"
	"practice-backend/internal/storage/inmem"
)

type Auth interface {
	Login(
		ctx context.Context,
		login string,
		password string,
	) (token string, err error)
	Register(
		ctx context.Context,
		user user.User,
	) (userID int, err error)
	IsAdmin(
		ctx context.Context,
		userID int,
	) (bool, error)
}

type HTTPHandlers struct {
	entryRepo   entry.EntryRepo
	userRepo    user.UserRepo
	authService Auth
}

func NewHTTPHandlers(
	entryRepo entry.EntryRepo,
	userRepo user.UserRepo,
	authService Auth,
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
	var registerDTO RegisterUserDTO

	if err := json.NewDecoder(r.Body).Decode(&registerDTO); err != nil {
		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.String(), http.StatusBadRequest)
		return
	}

	if err := registerDTO.Validate(); err != nil {
		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.String(), http.StatusBadRequest)
		return
	}

	user := user.NewUser(
		registerDTO.Login,
		registerDTO.Password,
		registerDTO.Name,
		registerDTO.Surname,
		registerDTO.Patronymic,
		registerDTO.Phone,
		registerDTO.Email,
		false,
	)

	userID, err := h.authService.Register(r.Context(), *user)
	if err != nil {
		errDTO := NewErrorDTO(err)
		if errors.Is(err, inmem.ErrUserAlreadyExist) {
			http.Error(w, errDTO.String(), http.StatusConflict)
			return
		}

		http.Error(w, errDTO.String(), http.StatusInternalServerError)
		return
	}

	resp := struct {
		ID int `json:"id"`
	}{
		ID: userID,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
