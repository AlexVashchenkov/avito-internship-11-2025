package service

import (
	"errors"

	api "github.com/AlexVashchenkov/avito-pr-reviewer-service/api"
	"github.com/AlexVashchenkov/avito-pr-reviewer-service/internal/repo"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserService struct {
	userRepo repo.UserRepository
	prRepo   repo.PullRequestRepository
}

func NewUserService(userRepo repo.UserRepository, prRepo repo.PullRequestRepository) *UserService {
	return &UserService{userRepo: userRepo, prRepo: prRepo}
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

	if newUser, ok := s.userRepo.UpdateUser(user); !ok {
		return nil, ErrPullRequestNotFound
	} else {
		return newUser, nil
	}
}
