package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	api "github.com/AlexVashchenkov/avito-pr-reviewer-service/api"
	"github.com/AlexVashchenkov/avito-pr-reviewer-service/internal/app"
	"github.com/AlexVashchenkov/avito-pr-reviewer-service/internal/repo/postgresql"
	"github.com/AlexVashchenkov/avito-pr-reviewer-service/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Unable to locate .env file")
	}

	connection, err := postgresql.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	userRepo := postgresql.NewUserRepository(connection)
	teamRepo := postgresql.NewTeamRepository(connection)
	prRepo := postgresql.NewPullRequestRepository(connection)

	prService := service.NewPullRequestService(prRepo, userRepo, teamRepo)
	userService := service.NewUserService(prRepo, userRepo, teamRepo)
	teamService := service.NewTeamService(teamRepo, userRepo)

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
