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

// Client wraps GitHub public REST API v3 calls.
// No authentication is required — only public activity is accessible.
type Client struct {
	httpClient *http.Client
}

// NewClient creates a GitHub API client using only public endpoints.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) newRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	return req, nil
}

// FetchRecentCommits returns commits from public repos for the given username
// in the last 7 days.
//
// Strategy:
//  1. Fetch PushEvents from /users/{username}/events/public (up to 100).
//  2. For each PushEvent within the time window, use the compare API
//     (/repos/{owner}/{repo}/compare/{before}...{head}) to get the real
//     commit count — the events payload does not include commits directly.
//
// Note: only activity from public repositories is visible without a token.
func (c *Client) FetchRecentCommits(username string, debug bool) ([]model.Commit, error) {
	cutoff := time.Now().UTC().AddDate(0, 0, -windowDays)
	url := fmt.Sprintf("%s/users/%s/events/public?per_page=100", apiBase, username)

	req, err := c.newRequest(url)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not connect to GitHub API")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("user %q not found on GitHub", username)
	}
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

		pushCommits, err := c.fetchCompareCommits(e.Repo.Name, e.Payload.Before, e.Payload.Head, e.CreatedAt)
		if err != nil {
			if debug {
				fmt.Printf("[debug] compare failed for %s (%s...%s): %v\n",
					e.Repo.Name, e.Payload.Before[:7], e.Payload.Head[:7], err)
			}
			// Fall back: count as 1 commit so we don't lose the push event
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
