package http

import (
	"encoding/json"
	"time"
)

type HelloWorldDTO struct {
	Title string `json:"title"`
}

type ErrorDTO struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

func NewErrorDTO(err error) ErrorDTO {
	return ErrorDTO{
		Message: err.Error(),
		Time:    time.Now(),
	}
}

func (e ErrorDTO) String() string {
	b, _ := json.MarshalIndent(e, "", "    ")

	return string(b)
}
