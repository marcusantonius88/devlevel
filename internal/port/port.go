// Package port defines the interfaces (ports) that the application core
// depends on. Adapters (e.g. the GitHub client) implement these interfaces,
// keeping the domain decoupled from any specific infrastructure.
package port

import "devlevel/internal/model"

// UserResolver resolves the identity of the authenticated user.
type UserResolver interface {
	GetAuthenticatedUser() (string, error)
}

// CommitFetcher fetches recent commits for a given username.
type CommitFetcher interface {
	FetchRecentCommits(username string, debug bool) ([]model.Commit, error)
}
