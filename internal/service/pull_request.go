package service

import (
	"context"
	"errors"
	"time"

	api "github.com/AlexVashchenkov/avito-pr-reviewer-service/api"
	"github.com/AlexVashchenkov/avito-pr-reviewer-service/internal/repo"
)

var (
	ErrPullRequestExists         = errors.New("pull-request already exists")
	ErrPullRequestNotFound       = errors.New("pull-request not found")
	ErrPullRequestOrUserNotFound = errors.New("pull-request or user not found")
	ErrPullRequestMergeError     = errors.New("unable to merge pull-request")
	ErrPullRequestAlreadyMerged  = errors.New("cannot reassign merged pull-request")
)

type PullRequestService struct {
	prRepo   repo.PullRequestRepository
	userRepo repo.UserRepository
	teamRepo repo.TeamRepository
}

func NewPullRequestService(prRepo repo.PullRequestRepository, userRepo repo.UserRepository, teamRepo repo.TeamRepository) *PullRequestService {
	return &PullRequestService{prRepo: prRepo, userRepo: userRepo, teamRepo: teamRepo}
}

func (s *PullRequestService) Create(ctx context.Context, req *api.PullRequestCreatePostReq) (*api.PullRequest, error) {
	if _, ok := s.prRepo.GetPullRequest(req.PullRequestID); ok {
		return nil, ErrPullRequestExists
	}

	author, ok := s.userRepo.GetUser(req.AuthorID)
	if !ok {
		return nil, ErrPullRequestOrUserNotFound
	}

	if _, ok := s.teamRepo.GetTeamByName(author.TeamName); !ok {
		return nil, ErrPullRequestOrUserNotFound
	}

	pr := &api.PullRequest{
		PullRequestID:     req.PullRequestID,
		PullRequestName:   req.PullRequestName,
		AuthorID:          req.AuthorID,
		Status:            api.PullRequestStatusOPEN,
		AssignedReviewers: s.prRepo.PickRandomReviewers(req.AuthorID, 2),
		CreatedAt:         api.NewOptNilDateTime(time.Now()),
	}

	s.prRepo.CreatePullRequest(pr)
	return pr, nil
}

func (s *PullRequestService) Merge(ctx context.Context, req *api.PullRequestMergePostReq) (*api.PullRequest, error) {
	pr, ok := s.prRepo.GetPullRequest(req.PullRequestID)
	if !ok {
		return nil, ErrPullRequestNotFound
	}

	if pr.Status == api.PullRequestStatusMERGED {
		return pr, nil
	}

	pr.SetStatus(api.PullRequestStatusMERGED)
	pr.SetMergedAt(api.NewOptNilDateTime(time.Now()))

	if newPr, ok := s.prRepo.UpdatePullRequest(pr); !ok {
		return nil, ErrPullRequestNotFound
	} else {
		return newPr, nil
	}
}

func (s *PullRequestService) Reassign(ctx context.Context, req *api.PullRequestReassignPostReq) (*api.PullRequest, error) {
	pr, ok := s.prRepo.GetPullRequest(req.PullRequestID)
	if !ok {
		return nil, ErrPullRequestNotFound
	}

	if pr.Status == api.PullRequestStatusMERGED {
		return nil, ErrPullRequestAlreadyMerged
	}
	return pr, nil
}
