package service

import (
	"context"
	"errors"
	"math/rand"
	"slices"
	"time"

	api "github.com/AlexVashchenkov/avito-pr-reviewer-service/api"
	"github.com/AlexVashchenkov/avito-pr-reviewer-service/internal/repo/contracts"
)

var (
	ErrPullRequestExists                = errors.New("pull-request already exists")
	ErrPullRequestNotFound              = errors.New("pull-request not found")
	ErrUserOrTeamNotFound               = errors.New("user or team not found")
	ErrPullRequestOrUserNotFound        = errors.New("pull-request or user not found")
	ErrPullRequestMergeError            = errors.New("unable to merge pull-request")
	ErrPullRequestAlreadyMerged         = errors.New("cannot reassign on merged PR")
	ErrPullRequestNoAvailableCandidates = errors.New("no active replacement candidate in team")
	ErrReviewerIsNotAssigned            = errors.New("reviewer is not assigned to this PR")
)

type PullRequestService struct {
	prRepo   contracts.PullRequestRepository
	userRepo contracts.UserRepository
	teamRepo contracts.TeamRepository
}

func NewPullRequestService(prRepo contracts.PullRequestRepository, userRepo contracts.UserRepository, teamRepo contracts.TeamRepository) *PullRequestService {
	return &PullRequestService{prRepo: prRepo, userRepo: userRepo, teamRepo: teamRepo}
}

func (s *PullRequestService) Create(ctx context.Context, req *api.PullRequestCreatePostReq) (*api.PullRequest, error) {
	if _, ok := s.prRepo.GetPullRequest(req.PullRequestID); ok {
		return nil, ErrPullRequestExists
	}

	author, ok := s.userRepo.GetUser(req.AuthorID)
	if !ok {
		return nil, ErrUserOrTeamNotFound
	}

	if _, ok := s.teamRepo.GetTeamByName(author.TeamName); !ok {
		return nil, ErrUserOrTeamNotFound
	}

	pr := &api.PullRequest{
		PullRequestID:     req.PullRequestID,
		PullRequestName:   req.PullRequestName,
		AuthorID:          req.AuthorID,
		Status:            api.PullRequestStatusOPEN,
		AssignedReviewers: s.userRepo.PickRandomReviewers(req.AuthorID, author.TeamName, 2),
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

func (s *PullRequestService) Reassign(ctx context.Context, req *api.PullRequestReassignPostReq) (*api.PullRequest, string, error) {
	pr, ok := s.prRepo.GetPullRequest(req.PullRequestID)
	if !ok {
		return nil, "", ErrPullRequestNotFound
	}

	if pr.Status == api.PullRequestStatusMERGED {
		return nil, "", ErrPullRequestAlreadyMerged
	}

	if isAssigned := slices.Contains(pr.AssignedReviewers, req.OldUserID); !isAssigned {
		return nil, "", ErrReviewerIsNotAssigned
	}

	oldReviewer, ok := s.userRepo.GetUser(req.OldUserID)
	if !ok {
		return nil, "", ErrPullRequestOrUserNotFound
	}

	team, ok := s.teamRepo.GetTeamByName(oldReviewer.TeamName)
	if !ok {
		return nil, "", ErrTeamNotFound
	}

	candidates := make([]string, 0)
	for _, member := range team.Members {
		if !member.IsActive || member.UserID == pr.AuthorID || member.UserID == req.OldUserID {
			continue
		}
		skip := slices.Contains(pr.AssignedReviewers, member.UserID)
		if !skip {
			candidates = append(candidates, member.UserID)
		}
	}

	if len(candidates) == 0 {
		return nil, "", ErrPullRequestNoAvailableCandidates
	}

	newReviewer := candidates[rand.Intn(len(candidates))]

	for i, rid := range pr.AssignedReviewers {
		if rid == req.OldUserID {
			pr.AssignedReviewers[i] = newReviewer
			break
		}
	}

	updated, ok := s.prRepo.UpdatePullRequest(pr)
	if !ok {
		return nil, "", ErrPullRequestNotFound
	}

	return updated, newReviewer, nil
}
