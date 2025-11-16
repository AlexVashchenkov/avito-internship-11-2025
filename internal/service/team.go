package service

import (
	"context"
	"errors"

	api "github.com/AlexVashchenkov/avito-pr-reviewer-service/api"
	"github.com/AlexVashchenkov/avito-pr-reviewer-service/internal/repo/contracts"
)

var (
	ErrTeamExists   = errors.New("team already exists")
	ErrTeamNotFound = errors.New("team not found")
)

type TeamService struct {
	teamRepo contracts.TeamRepository
}

func NewTeamService(teamRepo contracts.TeamRepository) *TeamService {
	return &TeamService{teamRepo: teamRepo}
}

func (s *TeamService) Create(ctx context.Context, team *api.Team) (*api.Team, error) {
	if _, ok := s.teamRepo.GetTeamByName(team.TeamName); ok {
		return nil, ErrTeamExists
	}

	s.teamRepo.CreateTeam(team)
	return team, nil
}

func (s *TeamService) GetByName(name string) (*api.Team, error) {
	team, ok := s.teamRepo.GetTeamByName(name)
	if !ok {
		return nil, ErrTeamNotFound
	}
	return team, nil
}
