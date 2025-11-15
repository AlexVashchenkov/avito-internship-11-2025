package repo

import api "github.com/AlexVashchenkov/avito-pr-reviewer-service/api"

type PullRequestRepository interface {
	CreatePullRequest(pr *api.PullRequest)
	UpdatePullRequest(pr *api.PullRequest) (*api.PullRequest, bool)
	GetPullRequest(id string) (*api.PullRequest, bool)
	GetPullRequestsByUser(userID string) []api.PullRequestShort
	PickRandomReviewers(authorID string, n int) []string
}
