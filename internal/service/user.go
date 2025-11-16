package service

import (
	"errors"

	api "github.com/AlexVashchenkov/avito-pr-reviewer-service/api"
	"github.com/AlexVashchenkov/avito-pr-reviewer-service/internal/repo/contracts"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService struct {
	userRepo contracts.UserRepository
	prRepo   contracts.PullRequestRepository
	teamRepo contracts.TeamRepository
}

func NewUserService(prRepo contracts.PullRequestRepository, userRepo contracts.UserRepository, teamRepo contracts.TeamRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		prRepo:   prRepo,
		teamRepo: teamRepo,
	}
}

func (s *UserService) GetUserReviews(id string) ([]api.PullRequestShort, error) {
	user, ok := s.userRepo.GetUser(id)
	if !ok {
		return nil, ErrUserNotFound
	}

	return s.prRepo.GetPullRequestsByUser(user.UserID), nil
}

func (s *UserService) SetUserActive(id string, isActive bool) (*api.User, error) {
	user, ok := s.userRepo.GetUser(id)
	if !ok {
		return nil, ErrUserNotFound
	}

	user.SetIsActive(isActive)

	if _, ok := s.teamRepo.UpdateTeamMember(user); !ok {
		return nil, ErrTeamNotFound
	}

	if _, ok := s.userRepo.UpdateUser(user); !ok {
		return nil, ErrUserNotFound
	}

	return user, nil
}
