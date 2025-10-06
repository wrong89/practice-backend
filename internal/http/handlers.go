package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"practice-backend/internal/models/entry"
	"practice-backend/internal/models/user"
	"practice-backend/internal/storage/inmem"
	"strconv"
	"time"
)

// todo: implement handler get user by login(pattern: /user/{login})
// todo: implement handler user is admin(pattern: /user/{login})

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

/*
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

func (h *HTTPHandlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
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

/*
pattern: /user/login
method:  POST
info:    JSON in HTTP request body

succeed:
  - status code: 200 OK
  - response body: JSON with token
failed:
  - status code: 400, 500
  - response body: JSON with error + time
*/

func (h *HTTPHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginDTO LoginUserDTO

	if err := json.NewDecoder(r.Body).Decode(&loginDTO); err != nil {
		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.String(), http.StatusBadRequest)
		return
	}

	if err := loginDTO.Validate(); err != nil {
		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.String(), http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(r.Context(), loginDTO.Login, loginDTO.Password)
	if err != nil {
		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.String(), http.StatusBadRequest)
		return
	}

	resp := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

/*
pattern: /entry
method:  POST
info:    JSON of created entry

succeed:
  - status code: 201 Created
  - response body: JSON with token
failed:
  - status code: 400, 500
  - response body: JSON with error + time
*/

func (h *HTTPHandlers) CreateEntryHandler(w http.ResponseWriter, r *http.Request) {
	var createEntryDTO CreateEntryDTO

	if err := json.NewDecoder(r.Body).Decode(&createEntryDTO); err != nil {
		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.String(), http.StatusBadRequest)
		return
	}

	if err := createEntryDTO.Validate(); err != nil {
		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.String(), http.StatusBadRequest)
		return
	}

	// skip the err, time is valid
	dateTime, _ := time.Parse(time.DateOnly, createEntryDTO.Date)

	entry, err := h.entryRepo.CreateEntry(
		r.Context(),
		createEntryDTO.Course,
		dateTime,
		createEntryDTO.UserID,
		createEntryDTO.PaymentMethod,
	)
	if err != nil {
		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.String(), http.StatusInternalServerError)
		return
	}

	resp := struct {
		Course        string `json:"course"`
		Date          string `json:"date"`
		UserID        int    `json:"user_id"`
		PaymentMethod string `json:"payment_method"`
		Status        string `json:"status"`
	}{
		Course:        entry.Course,
		Date:          createEntryDTO.Date,
		UserID:        entry.UserID,
		PaymentMethod: entry.PaymentMethod,
		Status:        entry.Status,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

/*
pattern: /entry?user_id=x
method:  GET
info:    query params

succeed:
  - status code: 200 OK
  - response body: JSON with entries
failed:
  - status code: 400, 500
  - response body: JSON with error + time
*/

func (h *HTTPHandlers) GetEntriesHandler(w http.ResponseWriter, r *http.Request) {
	entries, err := h.entryRepo.GetEntries(r.Context())
	if err != nil {
		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.String(), http.StatusInternalServerError)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	if len(userIDStr) > 0 {
		userID, err := strconv.Atoi(userIDStr)
		if err == nil {
			var filteredEntries []entry.Entry

			for _, entry := range entries {
				if entry.UserID == userID {
					filteredEntries = append(filteredEntries, entry)
				}
			}
			entries = filteredEntries
		}
	}

	resp := struct {
		Entries []entry.Entry `json:"entries"`
	}{
		Entries: entries,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

/*
pattern: /entry
method:  PATCH
info:    JSON with entry id and status

succeed:
  - status code: 200 OK
  - response body: JSON with entries
failed:
  - status code: 400, 500
  - response body: JSON with error + time
*/

func (h *HTTPHandlers) UpdateEntryHandler(w http.ResponseWriter, r *http.Request) {
	var updateEntryDTO UpdateEntryDTO

	if err := json.NewDecoder(r.Body).Decode(&updateEntryDTO); err != nil {
		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.String(), http.StatusBadRequest)
		return
	}

	if err := updateEntryDTO.Validate(); err != nil {
		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.String(), http.StatusBadRequest)
		return
	}

	entry, err := h.entryRepo.UpdateStatusEntry(r.Context(), updateEntryDTO.ID, updateEntryDTO.Status)
	if err != nil {
		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.String(), http.StatusBadRequest)
		return
	}

	resp := struct {
		Course        string `json:"course"`
		Date          string `json:"date"`
		UserID        int    `json:"user_id"`
		PaymentMethod string `json:"payment_method"`
		Status        string `json:"status"`
	}{
		Course:        entry.Course,
		Date:          entry.Date.Format(time.DateOnly),
		UserID:        entry.UserID,
		PaymentMethod: entry.PaymentMethod,
		Status:        entry.Status,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

/*
pattern: /entry/{user_id}
method:  GET
info:    in pattern

succeed:
  - status code: 200 OK
  - response body: JSON with entries
failed:
  - status code: 400, 404, 500
  - response body: JSON with error + time
*/

func (h *HTTPHandlers) GetEntriesByHandler(w http.ResponseWriter, r *http.Request) {
	var updateEntryDTO UpdateEntryDTO

	if err := json.NewDecoder(r.Body).Decode(&updateEntryDTO); err != nil {
		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.String(), http.StatusBadRequest)
		return
	}

	if err := updateEntryDTO.Validate(); err != nil {
		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.String(), http.StatusBadRequest)
		return
	}

	entry, err := h.entryRepo.UpdateStatusEntry(r.Context(), updateEntryDTO.ID, updateEntryDTO.Status)
	if err != nil {
		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.String(), http.StatusBadRequest)
		return
	}

	resp := struct {
		Course        string `json:"course"`
		Date          string `json:"date"`
		UserID        int    `json:"user_id"`
		PaymentMethod string `json:"payment_method"`
		Status        string `json:"status"`
	}{
		Course:        entry.Course,
		Date:          entry.Date.Format(time.DateOnly),
		UserID:        entry.UserID,
		PaymentMethod: entry.PaymentMethod,
		Status:        entry.Status,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

/*
pattern: /user/{user_id}
method:  GET
info:    in pattern

succeed:
  - status code: 200 OK
  - response body: JSON with isAdmin boolean
failed:
  - status code: 400, 404, 500
  - response body: JSON with error + time
*/

func (h *HTTPHandlers) UserIsAdminHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("user_id")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.String(), http.StatusBadRequest)
		return
	}

	isAdmin, err := h.authService.IsAdmin(r.Context(), userID)
	if err != nil {
		errDTO := NewErrorDTO(err)
		http.Error(w, errDTO.String(), http.StatusInternalServerError)
		return
	}

	resp := struct {
		IsAdmin bool `json:"is_admin"`
	}{
		IsAdmin: isAdmin,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
