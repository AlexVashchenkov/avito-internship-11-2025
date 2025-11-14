package main

import (
	"log"
	"net/http"

	api "github.com/AlexVashchenkov/avito-pr-reviewer-service/gen"
	"github.com/AlexVashchenkov/avito-pr-reviewer-service/internal/app"
	"github.com/AlexVashchenkov/avito-pr-reviewer-service/internal/repo"
	"github.com/AlexVashchenkov/avito-pr-reviewer-service/internal/service"
)

func main() {
	repo := repo.NewInMemoryRepository()

	prService := service.NewPullRequestService(repo)

	handler := app.NewHandler(prService)
	security := app.NewSecurityHandler()

	srv, err := api.NewServer(handler, security)
	if err != nil {
		log.Fatalf("failed to init server: %v", err)
	}

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
