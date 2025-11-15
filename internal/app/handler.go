package app

import (
	"context"
	"errors"

	api "github.com/AlexVashchenkov/avito-pr-reviewer-service/api"
	"github.com/AlexVashchenkov/avito-pr-reviewer-service/internal/service"
)

type Handler struct {
	prService   *service.PullRequestService
	userService *service.UserService
	teamService *service.TeamService
}

func NewHandler(prService *service.PullRequestService, userService *service.UserService, teamService *service.TeamService) *Handler {
	return &Handler{
		prService:   prService,
		userService: userService,
		teamService: teamService,
	}
}

func (h *Handler) PullRequestCreatePost(ctx context.Context, req *api.PullRequestCreatePostReq) (api.PullRequestCreatePostRes, error) {
	_, err := h.prService.Create(ctx, req)
	if errors.Is(err, service.ErrPullRequestExists) {
		return &api.PullRequestCreatePostConflict{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodePREXISTS,
				Message: err.Error(),
			},
		}, nil
	} else if errors.Is(err, service.ErrPullRequestOrUserNotFound) {
		return &api.PullRequestCreatePostNotFound{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeNOTFOUND,
				Message: err.Error(),
			},
		}, nil
	}
	return &api.PullRequestCreatePostCreated{}, nil
}

func (h *Handler) PullRequestMergePost(ctx context.Context, req *api.PullRequestMergePostReq) (api.PullRequestMergePostRes, error) {
	_, err := h.prService.Merge(ctx, req)
	if errors.Is(err, service.ErrPullRequestNotFound) {
		return &api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeNOTFOUND,
				Message: err.Error(),
			},
		}, nil
	}
	return &api.PullRequestMergePostOK{}, nil
}

func (h *Handler) PullRequestReassignPost(ctx context.Context, req *api.PullRequestReassignPostReq) (api.PullRequestReassignPostRes, error) {
	_, err := h.prService.Reassign(ctx, req)
	if errors.Is(err, service.ErrPullRequestNotFound) {
		return &api.PullRequestReassignPostNotFound{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeNOTFOUND,
				Message: err.Error(),
			},
		}, nil
	}
	return &api.PullRequestReassignPostOK{}, nil
}

func (h *Handler) TeamAddPost(ctx context.Context, req *api.Team) (api.TeamAddPostRes, error) {
	_, err := h.teamService.Create(ctx, req)
	if errors.Is(err, service.ErrTeamExists) {
		return &api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeTEAMEXISTS,
				Message: err.Error(),
			},
		}, nil
	}
	return &api.TeamAddPostCreated{}, nil
}

func (h *Handler) TeamGetGet(ctx context.Context, params api.TeamGetGetParams) (api.TeamGetGetRes, error) {
	team, err := h.teamService.GetByName(params.TeamName)
	if err != nil {
		return &api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeNOTFOUND,
				Message: err.Error(),
			},
		}, nil
	}
	return team, nil
}

func (h *Handler) UsersGetReviewGet(ctx context.Context, params api.UsersGetReviewGetParams) (*api.UsersGetReviewGetOK, error) {
	userPullRequests, err := h.userService.GetUserReviews(params.UserID)
	if errors.Is(err, service.ErrUserNotFound) {
		return &api.UsersGetReviewGetOK{
			UserID:       params.UserID,
			PullRequests: []api.PullRequestShort{},
		}, nil
	}
	return &api.UsersGetReviewGetOK{
		UserID:       params.UserID,
		PullRequests: userPullRequests,
	}, nil
}

func (h *Handler) UsersSetIsActivePost(ctx context.Context, req *api.UsersSetIsActivePostReq) (api.UsersSetIsActivePostRes, error) {
	user, err := h.userService.SetUserActive(req.UserID, req.IsActive)
	if errors.Is(err, service.ErrUserNotFound) {
		return &api.UsersSetIsActivePostNotFound{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeNOTFOUND,
				Message: err.Error(),
			},
		}, nil
	}
	return &api.UsersSetIsActivePostOK{
		User: api.NewOptUser(*user),
	}, nil
}
