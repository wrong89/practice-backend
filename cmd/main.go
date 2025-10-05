package main

import (
	"context"
	"fmt"
	"log"
	"practice-backend/internal/http"
	"practice-backend/internal/services/auth"
	"practice-backend/internal/storage/inmem"
)

func main() {
	storage := inmem.NewStorage()

	authService := auth.NewAuth(storage)

	authService.CreateAdminUser(context.TODO(), "Admin", "KorokNET")

	handlers := http.NewHTTPHandlers(storage, storage, authService)
	server := http.NewHTTPServer(*handlers)

	log.Printf("Starting server %s:%d\n", "localhost", 9091)

	if err := server.Start(); err != nil {
		fmt.Println("ERR", err)
	}
}
