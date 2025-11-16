package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/AlexVashchenkov/avito-pr-reviewer-service/api"
	"github.com/AlexVashchenkov/avito-pr-reviewer-service/internal/repo/memory"
	"github.com/AlexVashchenkov/avito-pr-reviewer-service/internal/service"
	"github.com/stretchr/testify/require"
)

func TestPullRequestService_CreateAlreadyExists(t *testing.T) {
	ctx := context.Background()
	prRepo := memory.InitPullRequestRepo()
	userRepo := memory.InitUserRepo()
	teamRepo := memory.InitTeamRepo()
	s := service.NewPullRequestService(prRepo, userRepo, teamRepo)

	req := &api.PullRequestCreatePostReq{
		PullRequestID:   "pr-test-1",
		PullRequestName: "Add new endpoint",
		AuthorID:        "u1",
	}

	pr, err := s.Create(ctx, req)
	require.NoError(t, err)
	require.Equal(t, "pr-test-1", pr.PullRequestID)
	require.Equal(t, api.PullRequestStatusOPEN, pr.Status)
	require.NotEmpty(t, pr.AssignedReviewers)

	req2 := &api.PullRequestCreatePostReq{
		PullRequestID:   "pr-test-1",
		PullRequestName: "Add another endpoint",
		AuthorID:        "u2",
	}

	_, err2 := s.Create(ctx, req2)
	require.Error(t, err2, service.ErrPullRequestExists)
}

func TestPullRequestService_AuthorNotFound(t *testing.T) {
	ctx := context.Background()
	prRepo := memory.InitPullRequestRepo()
	userRepo := memory.InitUserRepo()
	teamRepo := memory.InitTeamRepo()
	s := service.NewPullRequestService(prRepo, userRepo, teamRepo)

	req := &api.PullRequestCreatePostReq{
		PullRequestID:   "pr-test-1",
		PullRequestName: "Add new endpoint",
		AuthorID:        "u7",
	}

	_, err := s.Create(ctx, req)
	require.Equal(t, err, service.ErrUserOrTeamNotFound)
}

func TestPullRequestService_TeamNotFound(t *testing.T) {
	ctx := context.Background()
	prRepo := memory.InitPullRequestRepo()
	userRepo := memory.InitUserRepo()
	teamRepo := memory.InitTeamRepo()
	s := service.NewPullRequestService(prRepo, userRepo, teamRepo)

	req := &api.PullRequestCreatePostReq{
		PullRequestID:   "pr-test-1",
		PullRequestName: "Add new endpoint",
		AuthorID:        "u5",
	}

	_, err := s.Create(ctx, req)
	require.Equal(t, err, service.ErrUserOrTeamNotFound)
}

func TestPullRequestService_CreateAndMerge(t *testing.T) {
	ctx := context.Background()
	prRepo := memory.InitPullRequestRepo()
	userRepo := memory.InitUserRepo()
	teamRepo := memory.InitTeamRepo()
	s := service.NewPullRequestService(prRepo, userRepo, teamRepo)

	req := &api.PullRequestCreatePostReq{
		PullRequestID:   "pr-test-1",
		PullRequestName: "Add new endpoint",
		AuthorID:        "u1",
	}

	pr, err := s.Create(ctx, req)
	require.NoError(t, err)
	require.Equal(t, "pr-test-1", pr.PullRequestID)
	require.Equal(t, api.PullRequestStatusOPEN, pr.Status)
	require.NotEmpty(t, pr.AssignedReviewers)

	time.Sleep(10 * time.Millisecond)
	merged, err := s.Merge(ctx, &api.PullRequestMergePostReq{PullRequestID: pr.PullRequestID})
	require.NoError(t, err)
	require.Equal(t, api.PullRequestStatusMERGED, merged.Status)

	remerged, err := s.Merge(ctx, &api.PullRequestMergePostReq{PullRequestID: pr.PullRequestID})
	require.NoError(t, err)
	require.Equal(t, api.PullRequestStatusMERGED, remerged.Status)
	require.Equal(t, merged.GetMergedAt(), remerged.GetMergedAt())
}
