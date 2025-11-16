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
	if req.PullRequestID == "" {
		return &api.PullRequestCreatePostNotFound{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeNOTFOUND,
				Message: "pull_request_id is required",
			},
		}, nil
	}
	if req.PullRequestName == "" {
		return &api.PullRequestCreatePostNotFound{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeNOTFOUND,
				Message: "pull_request_name is required",
			},
		}, nil
	}
	if req.AuthorID == "" {
		return &api.PullRequestCreatePostNotFound{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeNOTFOUND,
				Message: "author_id is required",
			},
		}, nil
	}

	pr, err := h.prService.Create(ctx, req)
	if errors.Is(err, service.ErrPullRequestExists) {
		return &api.PullRequestCreatePostConflict{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodePREXISTS,
				Message: err.Error(),
			},
		}, nil
	} else if errors.Is(err, service.ErrUserOrTeamNotFound) {
		return &api.PullRequestCreatePostNotFound{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeNOTFOUND,
				Message: err.Error(),
			},
		}, nil
	}
	return &api.PullRequestCreatePostCreated{
		Pr: api.NewOptPullRequest(*pr),
	}, nil
}

func (h *Handler) PullRequestMergePost(ctx context.Context, req *api.PullRequestMergePostReq) (api.PullRequestMergePostRes, error) {
	if req.PullRequestID == "" {
		return &api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeNOTFOUND,
				Message: "pull_request_id is required",
			},
		}, nil
	}

	pr, err := h.prService.Merge(ctx, req)
	if errors.Is(err, service.ErrPullRequestNotFound) {
		return &api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeNOTFOUND,
				Message: err.Error(),
			},
		}, nil
	}
	return &api.PullRequestMergePostOK{
		Pr: api.NewOptPullRequest(*pr),
	}, nil
}

func (h *Handler) PullRequestReassignPost(ctx context.Context, req *api.PullRequestReassignPostReq) (api.PullRequestReassignPostRes, error) {
	if req.PullRequestID == "" {
		return &api.PullRequestReassignPostNotFound{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeNOTFOUND,
				Message: "pull_request_id is required",
			},
		}, nil
	}
	if req.OldUserID == "" {
		return &api.PullRequestReassignPostConflict{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeNOTFOUND,
				Message: "old_user_id is required",
			},
		}, nil
	}

	pr, newReviewer, err := h.prService.Reassign(ctx, req)
	if errors.Is(err, service.ErrPullRequestNotFound) {
		return &api.PullRequestReassignPostNotFound{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeNOTFOUND,
				Message: err.Error(),
			},
		}, nil
	} else if errors.Is(err, service.ErrPullRequestAlreadyMerged) {
		return &api.PullRequestReassignPostConflict{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodePRMERGED,
				Message: err.Error(),
			},
		}, nil
	} else if errors.Is(err, service.ErrReviewerIsNotAssigned) {
		return &api.PullRequestReassignPostConflict{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeNOTASSIGNED,
				Message: err.Error(),
			},
		}, nil
	} else if errors.Is(err, service.ErrPullRequestOrUserNotFound) {
		return &api.PullRequestReassignPostConflict{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeNOTFOUND,
				Message: err.Error(),
			},
		}, nil
	} else if errors.Is(err, service.ErrPullRequestNoAvailableCandidates) {
		return &api.PullRequestReassignPostConflict{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeNOCANDIDATE,
				Message: err.Error(),
			},
		}, nil
	}
	return &api.PullRequestReassignPostOK{
		Pr:         *pr,
		ReplacedBy: newReviewer,
	}, nil
}

func (h *Handler) TeamAddPost(ctx context.Context, req *api.Team) (api.TeamAddPostRes, error) {
	if req.TeamName == "" {
		return &api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeNOTFOUND,
				Message: "team_name is required",
			},
		}, nil
	}
	for i, member := range req.Members {
		if member.UserID == "" {
			return &api.ErrorResponse{
				Error: api.ErrorResponseError{
					Code:    api.ErrorResponseErrorCodeNOTFOUND,
					Message: "member user_id is required",
				},
			}, nil
		}
		if member.Username == "" {
			return &api.ErrorResponse{
				Error: api.ErrorResponseError{
					Code:    api.ErrorResponseErrorCodeNOTFOUND,
					Message: "member username is required",
				},
			}, nil
		}
		for j := i + 1; j < len(req.Members); j++ {
			if req.Members[j].UserID == member.UserID {
				return &api.ErrorResponse{
					Error: api.ErrorResponseError{
						Code:    api.ErrorResponseErrorCodeNOTFOUND,
						Message: "duplicate user_id in members",
					},
				}, nil
			}
		}
	}

	team, err := h.teamService.Create(ctx, req)
	if errors.Is(err, service.ErrTeamExists) {
		return &api.ErrorResponse{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeTEAMEXISTS,
				Message: err.Error(),
			},
		}, nil
	}
	return &api.TeamAddPostCreated{
		Team: api.NewOptTeam(*team),
	}, nil
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
	if req.UserID == "" {
		return &api.UsersSetIsActivePostNotFound{
			Error: api.ErrorResponseError{
				Code:    api.ErrorResponseErrorCodeNOTFOUND,
				Message: "user_id is required",
			},
		}, nil
	}

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
