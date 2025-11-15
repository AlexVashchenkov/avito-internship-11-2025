package service_test

/*
import (
	"context"
	"testing"

	"github.com/AlexVashchenkov/avito-pr-reviewer-service/api"
	"github.com/AlexVashchenkov/avito-pr-reviewer-service/internal/repo"
	"github.com/AlexVashchenkov/avito-pr-reviewer-service/internal/service"
	"github.com/stretchr/testify/require"
)

func TestTeamService_Create(t *testing.T) {
	ctx := context.Background()
	r := repo.InitRepo()
	s := service.NewTeamService(r)

	team := &api.Team{
		TeamName: "Product",
		Members: []api.TeamMember{
			{
				UserID:   "u5",
				Username: "Eugene",
				IsActive: true,
			},
		},
	}

	teamRes, err := s.Create(ctx, team)
	require.NoError(t, err)
	require.Equal(t, teamRes.TeamName, team.TeamName)
	require.Equal(t, teamRes.Members, team.Members)
}

func TestTeamService_CreateAlreadyExists(t *testing.T) {
	ctx := context.Background()
	r := repo.InitRepo()
	s := service.NewTeamService(r)

	team := &api.Team{
		TeamName: "frontend",
		Members: []api.TeamMember{
			{
				UserID:   "u5",
				Username: "Eugene",
				IsActive: true,
			},
		},
	}

	_, err := s.Create(ctx, team)
	require.Error(t, err)
	require.Equal(t, err, service.ErrTeamExists)
}
*/
