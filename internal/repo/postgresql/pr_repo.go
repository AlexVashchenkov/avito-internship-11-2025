package postgresql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"

	api "github.com/AlexVashchenkov/avito-pr-reviewer-service/api"
)

type PullRequestRepository struct {
	db *sql.DB
}

func NewPullRequestRepository(db *sql.DB) *PullRequestRepository {
	return &PullRequestRepository{db: db}
}

func (r *PullRequestRepository) CreatePullRequest(pr *api.PullRequest) {
	reviewers, _ := json.Marshal(pr.AssignedReviewers)

	_, err := r.db.Exec(`
		INSERT INTO pull_requests (
			pull_request_id, pull_request_name, author_id,
			status_id, assigned_reviewers, created_at, merged_at
		)
		VALUES ($1, $2, $3,
		        (SELECT id FROM pr_statuses WHERE status=$4),
		        $5::jsonb, $6, $7)
	`,
		pr.PullRequestID,
		pr.PullRequestName,
		pr.AuthorID,
		pr.Status,
		reviewers,
		pr.CreatedAt.Value,
		pr.MergedAt.Value,
	)

	if err != nil {
		slog.Error("failed to create pull request",
			"pull_request_id", pr.PullRequestID,
			"author_id", pr.AuthorID,
			"error", err.Error(),
		)
		return
	}

	slog.Info("pull request created",
		"pull_request_id", pr.PullRequestID,
		"author_id", pr.AuthorID,
	)
}

func (r *PullRequestRepository) GetPullRequest(id string) (*api.PullRequest, bool) {
	row := r.db.QueryRow(`
		SELECT pr.pull_request_id, pr.pull_request_name,
		       pr.author_id, ps.status,
		       pr.assigned_reviewers, pr.created_at, pr.merged_at
		FROM pull_requests pr
		JOIN pr_statuses ps ON ps.id = pr.status_id
		WHERE pr.pull_request_id = $1
	`, id)

	var reviewersJSON []byte
	var pr api.PullRequest
	var status string

	err := row.Scan(
		&pr.PullRequestID,
		&pr.PullRequestName,
		&pr.AuthorID,
		&status,
		&reviewersJSON,
		&pr.CreatedAt.Value,
		&pr.MergedAt.Value,
	)
	if errors.Is(err, sql.ErrNoRows) {
		slog.Info("pull request not found", "pull_request_id", id)
		return nil, false
	}
	if err != nil {
		slog.Error("failed to get pull request",
			"pull_request_id", id,
			"error", err.Error(),
		)
		return nil, false
	}

	_ = json.Unmarshal(reviewersJSON, &pr.AssignedReviewers)
	pr.Status = api.PullRequestStatus(status)

	return &pr, true
}

func (r *PullRequestRepository) UpdatePullRequest(pr *api.PullRequest) (*api.PullRequest, bool) {
	reviewers, _ := json.Marshal(pr.AssignedReviewers)

	_, err := r.db.Exec(`
		UPDATE pull_requests
		SET pull_request_name=$2,
		    status_id = (SELECT id FROM pr_statuses WHERE status=$3),
		    assigned_reviewers=$4::jsonb,
		    merged_at=$5
		WHERE pull_request_id=$1
	`,
		pr.PullRequestID,
		pr.PullRequestName,
		pr.Status,
		reviewers,
		pr.MergedAt.Value,
	)

	if err != nil {
		slog.Error("failed to update pull request",
			"pull_request_id", pr.PullRequestID,
			"error", err.Error(),
		)
		return nil, false
	}

	slog.Info("pull request updated",
		"pull_request_id", pr.PullRequestID,
	)

	return pr, true
}

func (r *PullRequestRepository) GetPullRequestsByUser(userID string) []api.PullRequestShort {
	rows, err := r.db.Query(`
		SELECT pr.pull_request_id, pr.pull_request_name,
			   pr.author_id, ps.status
		FROM pull_requests pr
		JOIN pr_statuses ps ON ps.id = pr.status_id
		WHERE pr.assigned_reviewers @> jsonb_build_array($1::text)
	`, userID)

	if err != nil {
		slog.Error("failed to get PR list for user", "user_id", userID, "error", err.Error())
		return nil
	}
	defer rows.Close()

	var res []api.PullRequestShort
	for rows.Next() {
		var out api.PullRequestShort
		if scanErr := rows.Scan(&out.PullRequestID, &out.PullRequestName, &out.AuthorID, &out.Status); scanErr != nil {
			slog.Error("failed to scan PR row", "error", scanErr.Error())
			continue
		}
		res = append(res, out)
	}

	if err := rows.Err(); err != nil {
		slog.Error("error iterating PR rows", "user_id", userID, "error", err.Error())
		return nil
	}

	slog.Info("loaded user pull requests", "user_id", userID, "count", len(res))
	return res
}
