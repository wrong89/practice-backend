package main

import (
	"fmt"
	"practice-backend/internal/http"
	"practice-backend/internal/storage/inmem"
)

func main() {

	storage := inmem.NewStorage()

	handlers := http.NewHTTPHandlers(storage, storage)
	server := http.NewHTTPServer(*handlers)

	if err := server.Start(); err != nil {
		fmt.Println("ERR", err)
	}
}
