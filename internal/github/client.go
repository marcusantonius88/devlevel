package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"devlevel/internal/model"
)

const (
	apiBase    = "https://api.github.com"
	windowDays = 7
)

// Client wraps GitHub REST API v3 calls.
type Client struct {
	token      string
	httpClient *http.Client
}

// NewClient creates an authenticated GitHub API client.
func NewClient(token string) *Client {
	return &Client{
		token:      token,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) newRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	return req, nil
}

// GetAuthenticatedUser resolves the GitHub login for the supplied token.
func (c *Client) GetAuthenticatedUser() (string, error) {
	req, err := c.newRequest(apiBase + "/user")
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not connect to GitHub API")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("GITHUB_TOKEN is missing or invalid")
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API error: status %d", resp.StatusCode)
	}

	var user struct {
		Login string `json:"login"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", err
	}

	return user.Login, nil
}

// FetchRecentCommits returns commits authored by username in the last 7 days.
//
// Strategy:
//  1. Fetch PushEvents from /users/{username}/events (up to 100).
//  2. For each PushEvent within the time window, use the compare API
//     (/repos/{owner}/{repo}/compare/{before}...{head}) to get the real
//     commit list — this works even for private repos when the token has
//     access, and avoids the empty-commits-array issue in the events payload.
func (c *Client) FetchRecentCommits(username string, debug bool) ([]model.Commit, error) {
	cutoff := time.Now().UTC().AddDate(0, 0, -windowDays)
	url := fmt.Sprintf("%s/users/%s/events?per_page=100", apiBase, username)

	req, err := c.newRequest(url)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not connect to GitHub API")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: status %d", resp.StatusCode)
	}

	var events []struct {
		Type      string    `json:"type"`
		CreatedAt time.Time `json:"created_at"`
		Repo      struct {
			Name string `json:"name"` // "owner/repo"
		} `json:"repo"`
		Payload struct {
			Before string `json:"before"`
			Head   string `json:"head"`
		} `json:"payload"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, err
	}

	var commits []model.Commit
	for _, e := range events {
		if e.Type != "PushEvent" {
			continue
		}
		if e.CreatedAt.Before(cutoff) {
			continue
		}
		if e.Payload.Before == "" || e.Payload.Head == "" {
			continue
		}

		// Use the compare API to get the actual commits in this push.
		pushCommits, err := c.fetchCompareCommits(e.Repo.Name, e.Payload.Before, e.Payload.Head, e.CreatedAt)
		if err != nil {
			if debug {
				fmt.Printf("[debug] compare failed for %s (%s...%s): %v\n",
					e.Repo.Name, e.Payload.Before[:7], e.Payload.Head[:7], err)
			}
			// Fall back: count as 1 commit for this push so we don't lose the event
			commits = append(commits, model.Commit{SHA: e.Payload.Head, Date: e.CreatedAt})
			continue
		}

		if debug {
			fmt.Printf("[debug] PushEvent at %s — repo: %s — %d commit(s)\n",
				e.CreatedAt.Format("2006-01-02"), e.Repo.Name, len(pushCommits))
		}
		commits = append(commits, pushCommits...)
	}

	return commits, nil
}

// fetchCompareCommits calls the compare API and returns the commits between
// base and head, all stamped with the push timestamp.
func (c *Client) fetchCompareCommits(repoFullName, base, head string, pushTime time.Time) ([]model.Commit, error) {
	url := fmt.Sprintf("%s/repos/%s/compare/%s...%s", apiBase, repoFullName, base, head)

	req, err := c.newRequest(url)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("compare API status %d", resp.StatusCode)
	}

	var result struct {
		Commits []struct {
			SHA string `json:"sha"`
		} `json:"commits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	commits := make([]model.Commit, len(result.Commits))
	for i, rc := range result.Commits {
		commits[i] = model.Commit{SHA: rc.SHA, Date: pushTime}
	}
	return commits, nil
}
