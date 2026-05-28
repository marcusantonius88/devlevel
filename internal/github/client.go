package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"devlevel/internal/model"
	"devlevel/internal/port"
)

const (
	apiBase    = "https://api.github.com"
	windowDays = 30
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
// in the last windowDays days.
//
// Strategy (request-efficient):
//  1. Fetch PushEvents to discover which repos had activity (1 request).
//  2. Deduplicate repos touched within the window.
//  3. For each unique repo, fetch commits filtered by author + since date
//     (1 request per repo) — far fewer requests than 1 per PushEvent.
//
// Note: only activity from public repositories is visible without a token.
func (c *Client) FetchRecentCommits(username string, debug bool) ([]model.Commit, error) {
	since := time.Now().UTC().AddDate(0, 0, -windowDays)

	// Step 1 — discover active repos from PushEvents (1 request).
	activeRepos, err := c.fetchActiveRepos(username, since, debug)
	if err != nil {
		return nil, err
	}

	if debug {
		fmt.Printf("[debug] active repos in window: %v\n", activeRepos)
	}

	// Step 2 — fetch commits per repo (1 request each).
	var commits []model.Commit
	for _, repo := range activeRepos {
		repoCommits, err := c.fetchRepoCommits(repo, username, since, debug)
		if err != nil {
			if debug {
				fmt.Printf("[debug] skipping repo %s: %v\n", repo, err)
			}
			continue
		}
		commits = append(commits, repoCommits...)
	}

	return commits, nil
}

// fetchActiveRepos returns the unique list of public repos the user pushed to
// within the time window. Uses a single events/public request.
func (c *Client) fetchActiveRepos(username string, since time.Time, debug bool) ([]string, error) {
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
	if resp.StatusCode == http.StatusForbidden {
		return nil, port.ErrRateLimit
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
	}

	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var repos []string
	for _, e := range events {
		if e.Type != "PushEvent" {
			continue
		}
		if e.CreatedAt.Before(since) {
			continue
		}
		if !seen[e.Repo.Name] {
			seen[e.Repo.Name] = true
			repos = append(repos, e.Repo.Name)
		}
	}

	return repos, nil
}

// fetchRepoCommits fetches commits by the given author in a repo since a date.
// Uses /repos/{owner}/{repo}/commits?author=...&since=... (1 request per repo).
func (c *Client) fetchRepoCommits(repoFullName, author string, since time.Time, debug bool) ([]model.Commit, error) {
	url := fmt.Sprintf("%s/repos/%s/commits?author=%s&since=%s&per_page=100",
		apiBase, repoFullName, author, since.Format(time.RFC3339))

	req, err := c.newRequest(url)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("rate limit or access denied")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("commits API status %d", resp.StatusCode)
	}

	var result []struct {
		SHA    string `json:"sha"`
		Commit struct {
			Author struct {
				Date time.Time `json:"date"`
			} `json:"author"`
		} `json:"commit"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	commits := make([]model.Commit, len(result))
	for i, rc := range result {
		commits[i] = model.Commit{
			SHA:  rc.SHA,
			Date: rc.Commit.Author.Date,
		}
	}

	if debug {
		fmt.Printf("[debug] repo: %s — %d commit(s) in window\n", repoFullName, len(commits))
	}

	return commits, nil
}
