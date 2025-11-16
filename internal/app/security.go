package app

import (
	"context"
	"errors"

	api "github.com/AlexVashchenkov/avito-pr-reviewer-service/api"
)

type SecurityHandler struct {
}

func NewSecurityHandler() *SecurityHandler {
	return &SecurityHandler{}
}

func (s *SecurityHandler) HandleAdminToken(ctx context.Context, op api.OperationName, t api.AdminToken) (context.Context, error) {
	if t.Token == "" {
		return ctx, errors.New("missing admin token")
	}
	return ctx, nil
}

func (s *SecurityHandler) HandleUserToken(ctx context.Context, op api.OperationName, t api.UserToken) (context.Context, error) {
	if t.Token == "" {
		return ctx, errors.New("missing user token")
	}
	return ctx, nil
}
