// Package port defines the interfaces (ports) that the application core
// depends on. Adapters (e.g. the GitHub client) implement these interfaces,
// keeping the domain decoupled from any specific infrastructure.
package port

import (
	"errors"

	"devlevel/internal/model"
)

// ErrRateLimit is returned when the GitHub API rate limit has been reached.
var ErrRateLimit = errors.New("rate_limit")

// FetchResult holds the commits retrieved and metadata about the fetch.
type FetchResult struct {
	Commits        []model.Commit
	TotalRepos     int // repos discovered in the activity window
	SkippedRepos   int // repos skipped due to timeout or error
}

// CommitFetcher fetches recent commits for a given username.
type CommitFetcher interface {
	FetchRecentCommits(username string, debug bool) (FetchResult, error)
}
