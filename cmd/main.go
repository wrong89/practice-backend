package main

import (
	"context"
	"fmt"
	"practice-backend/internal/http"
	"practice-backend/internal/services/auth"
	"practice-backend/internal/storage/inmem"
)

func main() {
	storage := inmem.NewStorage()
	storage.CreateAdmin(context.TODO(), "Admin", "KorokNET")

	authService := auth.NewAuth(storage)

	handlers := http.NewHTTPHandlers(storage, storage, *authService)
	server := http.NewHTTPServer(*handlers)

	if err := server.Start(); err != nil {
		fmt.Println("ERR", err)
	}
}
