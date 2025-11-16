package main

import (
	"log"
	"net/http"

	api "github.com/AlexVashchenkov/avito-pr-reviewer-service/api"
	"github.com/AlexVashchenkov/avito-pr-reviewer-service/internal/app"
	"github.com/AlexVashchenkov/avito-pr-reviewer-service/internal/repo/memory"
	"github.com/AlexVashchenkov/avito-pr-reviewer-service/internal/service"
)

func main() {
	prRepo := memory.InitPullRequestRepo()
	userRepo := memory.InitUserRepo()
	teamRepo := memory.InitTeamRepo()

	prService := service.NewPullRequestService(prRepo, userRepo, teamRepo)
	userService := service.NewUserService(prRepo, userRepo, teamRepo)
	teamService := service.NewTeamService(teamRepo)

	handler := app.NewHandler(prService, userService, teamService)
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
