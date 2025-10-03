package http

import (
	"encoding/json"
	"net/http"
	"practice-backend/internal/models/entry"
	"practice-backend/internal/models/user"
)

type HTTPHandlers struct {
	entryRepo entry.EntryRepo
	userRepo  user.UserRepo
}

func NewHTTPHandlers(entryRepo entry.EntryRepo, userRepo user.UserRepo) *HTTPHandlers {
	return &HTTPHandlers{
		entryRepo: entryRepo,
		userRepo:  userRepo,
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
pattern: /books
method:  POST
info:    JSON in HTTP request body

succeed:
  - status code: 201 Created
  - response body: JSON of created book

failed:
  - status code: 400, 409(Conflict), 500
  - response body: JSON with error + time
*/
