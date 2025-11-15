package repo

import api "github.com/AlexVashchenkov/avito-pr-reviewer-service/api"

type UserRepo interface {
	CreateUser(user *api.User)
	UpdateUser(user *api.User) (*api.User, bool)
	GetUser(id string) (*api.User, bool)
}
