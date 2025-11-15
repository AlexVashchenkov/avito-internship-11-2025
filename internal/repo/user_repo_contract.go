package repo

import api "github.com/AlexVashchenkov/avito-pr-reviewer-service/api"

type UserRepository interface {
	CreateUser(user *api.User)
	GetUser(id string) (*api.User, bool)
	UpdateUser(user *api.User) (*api.User, bool)
}
