package app

import (
	"context"

	api "github.com/AlexVashchenkov/avito-pr-reviewer-service/gen"
	"github.com/AlexVashchenkov/avito-pr-reviewer-service/internal/service"
)

type Handler struct {
	prService *service.PullRequestService
}

func NewHandler(prService *service.PullRequestService) *Handler {
	return &Handler{
		prService: prService,
	}
}

func (h *Handler) PullRequestCreatePost(ctx context.Context, req *api.PullRequestCreatePostReq) (api.PullRequestCreatePostRes, error) {
	return &api.PullRequestCreatePostCreated{}, nil
}

func (h *Handler) PullRequestMergePost(ctx context.Context, req *api.PullRequestMergePostReq) (api.PullRequestMergePostRes, error) {
	return &api.PullRequestMergePostOK{}, nil
}

func (h *Handler) PullRequestReassignPost(ctx context.Context, req *api.PullRequestReassignPostReq) (api.PullRequestReassignPostRes, error) {
	return &api.PullRequestReassignPostOK{}, nil
}

func (h *Handler) TeamAddPost(ctx context.Context, req *api.Team) (api.TeamAddPostRes, error) {
	return &api.TeamAddPostCreated{}, nil
}

func (h *Handler) TeamGetGet(ctx context.Context, params api.TeamGetGetParams) (api.TeamGetGetRes, error) {
	return &api.Team{}, nil
}

func (h *Handler) UsersGetReviewGet(ctx context.Context, params api.UsersGetReviewGetParams) (*api.UsersGetReviewGetOK, error) {
	return &api.UsersGetReviewGetOK{}, nil
}

func (h *Handler) UsersSetIsActivePost(ctx context.Context, req *api.UsersSetIsActivePostReq) (api.UsersSetIsActivePostRes, error) {
	return &api.UsersSetIsActivePostOK{}, nil
}
