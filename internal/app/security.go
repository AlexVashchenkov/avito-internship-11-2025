package app

import (
	"context"

	api "github.com/AlexVashchenkov/avito-pr-reviewer-service/gen"
)

type SecurityHandler struct {
}

func NewSecurityHandler() *SecurityHandler {
	return &SecurityHandler{}
}

func (s *SecurityHandler) HandleAdminToken(ctx context.Context, operationName api.OperationName, t api.AdminToken) (context.Context, error) {
	return ctx, nil
}

func (s *SecurityHandler) HandleUserToken(ctx context.Context, operationName api.OperationName, t api.UserToken) (context.Context, error) {
	return ctx, nil
}
