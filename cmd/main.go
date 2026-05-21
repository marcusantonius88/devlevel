package main

import (
	"fmt"
	"os"

	"devlevel/internal/env"
	"devlevel/internal/gamification"
	githubadapter "devlevel/internal/github"
	"devlevel/internal/model"
	"devlevel/internal/port"
	"devlevel/internal/ui"
)

func main() {
	if err := env.Load(".env"); err != nil {
		fmt.Fprintln(os.Stderr, "Warning: could not read .env file:", err)
	}

	debug := len(os.Args) > 1 && os.Args[1] == "--debug"

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "Error: GITHUB_TOKEN is not set.")
		fmt.Fprintln(os.Stderr, "Tip: add it to the .env file as GITHUB_TOKEN=ghp_...")
		os.Exit(1)
	}

	// Wire the GitHub adapter to the port interfaces.
	// main.go only knows about port.UserResolver and port.CommitFetcher —
	// swapping the adapter (e.g. for GitLab) requires no changes here.
	adapter := githubadapter.NewClient(token)
	run(adapter, adapter, debug)
}

// run contains the application flow and depends only on port interfaces,
// making it straightforward to test with mock adapters.
func run(resolver port.UserResolver, fetcher port.CommitFetcher, debug bool) {
	username, err := resolver.GetAuthenticatedUser()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	commits, err := fetcher.FetchRecentCommits(username, debug)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	stats := buildStats(username, commits)

	ui.Render(
		stats,
		gamification.RankTitle(stats.Level),
		gamification.LevelProgress(stats.XP),
		gamification.XPToNextLevel(stats.XP),
	)
}

// buildStats computes all metrics from raw commits and returns a Stats value.
func buildStats(username string, commits []model.Commit) model.Stats {
	stats := model.Stats{
		Username:     username,
		XP:           gamification.CalculateXP(commits),
		Streak:       gamification.CalculateStreak(commits),
		CommitCount:  len(commits),
		DailyGoalMet: gamification.DailyGoalMet(commits),
	}
	stats.Level = gamification.CalculateLevel(stats.XP)
	return stats
}
