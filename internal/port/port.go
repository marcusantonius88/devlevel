// Package port defines the interfaces (ports) that the application core
// depends on. Adapters (e.g. the GitHub client) implement these interfaces,
// keeping the domain decoupled from any specific infrastructure.
package port

import "devlevel/internal/model"

// CommitFetcher fetches recent commits for a given username.
type CommitFetcher interface {
	FetchRecentCommits(username string, debug bool) ([]model.Commit, error)
}
