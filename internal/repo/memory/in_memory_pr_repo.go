package memory

import (
	api "github.com/AlexVashchenkov/avito-pr-reviewer-service/api"
)

type InMemoryPullRequestRepository struct {
	prs map[string]*api.PullRequest
}

func NewInMemoryPullRequestRepository() *InMemoryPullRequestRepository {
	return &InMemoryPullRequestRepository{
		prs: make(map[string]*api.PullRequest),
	}
}

func (r *InMemoryPullRequestRepository) CreatePullRequest(pr *api.PullRequest) {
	r.prs[pr.PullRequestID] = pr
}

func (r *InMemoryPullRequestRepository) GetPullRequest(id string) (*api.PullRequest, bool) {
	pr, ok := r.prs[id]
	return pr, ok
}

func (r *InMemoryPullRequestRepository) UpdatePullRequest(pr *api.PullRequest) (*api.PullRequest, bool) {
	_, ok := r.prs[pr.PullRequestID]
	if !ok {
		return nil, false
	}
	r.prs[pr.PullRequestID] = pr
	return pr, true
}

func (r *InMemoryPullRequestRepository) GetPullRequestsByUser(userID string) []api.PullRequestShort {
	result := make([]api.PullRequestShort, 0)
	for _, pr := range r.prs {
		for _, reviewerID := range pr.AssignedReviewers {
			if reviewerID == userID {
				shortPr := api.PullRequestShort{
					PullRequestID:   pr.PullRequestID,
					PullRequestName: pr.PullRequestName,
					AuthorID:        pr.AuthorID,
					Status:          api.PullRequestShortStatus(pr.Status),
				}
				result = append(result, shortPr)
				break
			}
		}
	}
	return result
}

func InitPullRequestRepo() *InMemoryPullRequestRepository {
	repo := NewInMemoryPullRequestRepository()

	return repo
}
