package repo

import (
	api "github.com/AlexVashchenkov/avito-pr-reviewer-service/gen"
)

type InMemoryRepository struct {
	prs   map[string]*api.PullRequest
	users map[string]*api.User
	teams map[string]*api.Team
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		prs:   make(map[string]*api.PullRequest),
		users: make(map[string]*api.User),
		teams: make(map[string]*api.Team),
	}
}

func (r *InMemoryRepository) Create(pr *api.PullRequest) {
	r.prs[pr.PullRequestID] = pr
}

func (r *InMemoryRepository) CheckExists(id string) bool {
	_, ok := r.prs[id]
	return ok
}

func (r *InMemoryRepository) PickRandomReviewers(authorID string, n int) []string {
	result := make([]string, 0, n)
	for userID := range r.users {
		if userID != authorID {
			result = append(result, userID)
		}
		if len(result) == n {
			break
		}
	}
	return result
}
