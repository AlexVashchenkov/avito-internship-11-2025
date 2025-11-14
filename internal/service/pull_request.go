package service

import (
	"context"
	"errors"
	"time"

	api "github.com/AlexVashchenkov/avito-pr-reviewer-service/gen"
	"github.com/AlexVashchenkov/avito-pr-reviewer-service/internal/repo"
)

type PullRequestService struct {
	repo *repo.InMemoryRepository
}

func NewPullRequestService(repo *repo.InMemoryRepository) *PullRequestService {
	return &PullRequestService{repo: repo}
}

func (s *PullRequestService) CreatePullRequest(ctx context.Context, req *api.PullRequestCreatePostReq) (*api.PullRequest, error) {
	if s.repo.CheckExists(req.PullRequestID) {
		return nil, errors.New("pull request already exists")
	}

	pr := &api.PullRequest{
		PullRequestID:     req.PullRequestID,
		PullRequestName:   req.PullRequestName,
		AuthorID:          req.AuthorID,
		Status:            api.PullRequestStatusOPEN,
		AssignedReviewers: s.repo.PickRandomReviewers(req.AuthorID, 2),
		CreatedAt:         api.NewOptNilDateTime(time.Now()),
	}

	s.repo.Create(pr)
	return pr, nil
}
