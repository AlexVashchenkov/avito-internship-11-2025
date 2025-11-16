package memory

import (
	api "github.com/AlexVashchenkov/avito-pr-reviewer-service/api"
)

type InMemoryUserRepository struct {
	users map[string]*api.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*api.User),
	}
}

func (r *InMemoryUserRepository) CreateUser(user *api.User) {
	r.users[user.UserID] = user
}

func (r *InMemoryUserRepository) GetUser(id string) (*api.User, bool) {
	user, ok := r.users[id]
	return user, ok
}

func (r *InMemoryUserRepository) UpdateUser(user *api.User) (*api.User, bool) {
	_, ok := r.users[user.UserID]
	if !ok {
		return nil, false
	}
	r.users[user.UserID] = user

	return user, ok
}

func (r *InMemoryUserRepository) PickRandomReviewers(authorID string, name string, n int) []string {
	result := make([]string, 0, n)
	for userID, user := range r.users {
		if userID != authorID && user.TeamName == name && user.IsActive == true {
			result = append(result, user.UserID)
		}
		if len(result) == n {
			break
		}
	}
	return result
}

func InitUserRepo() *InMemoryUserRepository {
	repo := NewInMemoryUserRepository()

	users := []*api.User{
		{UserID: "u1", Username: "Alice", TeamName: "backend", IsActive: true},
		{UserID: "u2", Username: "Bob", TeamName: "backend", IsActive: true},
		{UserID: "u3", Username: "Charlie", TeamName: "backend", IsActive: false},
		{UserID: "u4", Username: "Dana", TeamName: "frontend", IsActive: true},
		{UserID: "u5", Username: "Eugene", TeamName: "qa", IsActive: true},
	}

	for _, u := range users {
		repo.CreateUser(u)
	}

	return repo
}
