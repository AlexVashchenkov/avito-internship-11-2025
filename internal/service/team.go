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
	userRepo contracts.UserRepository
}

func NewTeamService(teamRepo contracts.TeamRepository, userRepo contracts.UserRepository) *TeamService {
	return &TeamService{teamRepo: teamRepo}
}

func NewTeamServiceWithUserRepo(teamRepo contracts.TeamRepository, userRepo contracts.UserRepository) *TeamService {
	return &TeamService{teamRepo: teamRepo, userRepo: userRepo}
}

func (s *TeamService) Create(ctx context.Context, team *api.Team) (*api.Team, error) {
	if _, ok := s.teamRepo.GetTeamByName(team.TeamName); ok {
		return nil, ErrTeamExists
	}

	if repoWithTx, ok := s.teamRepo.(interface {
		CreateTeamWithMembers(team *api.Team) error
	}); ok {
		if err := repoWithTx.CreateTeamWithMembers(team); err != nil {
			return nil, err
		}
	} else {
		s.teamRepo.CreateTeam(team)
		if s.userRepo != nil {
			for _, member := range team.Members {
				user := &api.User{
					UserID:   member.UserID,
					Username: member.Username,
					TeamName: team.TeamName,
					IsActive: member.IsActive,
				}

				if _, ok := s.userRepo.GetUser(member.UserID); ok {
					s.userRepo.UpdateUser(user)
				} else {
					s.userRepo.CreateUser(user)
				}
			}
		}
	}

	return team, nil
}

func (s *TeamService) GetByName(name string) (*api.Team, error) {
	team, ok := s.teamRepo.GetTeamByName(name)
	if !ok {
		return nil, ErrTeamNotFound
	}
	return team, nil
}
