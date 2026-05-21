package main

import (
	"flag"
	"fmt"
	"os"

	"devlevel/internal/gamification"
	githubadapter "devlevel/internal/github"
	"devlevel/internal/model"
	"devlevel/internal/port"
	"devlevel/internal/ui"
)

func main() {
	username := flag.String("user", "", "GitHub username (e.g. --user marcusantonius88)")
	debug := flag.Bool("debug", false, "Print debug information about API calls")
	flag.Parse()

	if *username == "" {
		fmt.Fprintln(os.Stderr, "Error: --user is required.")
		fmt.Fprintln(os.Stderr, "Usage: devlevel --user <github-username>")
		os.Exit(1)
	}

	adapter := githubadapter.NewClient()
	run(*username, adapter, *debug)
}

// run contains the application flow and depends only on port interfaces,
// making it straightforward to test with mock adapters.
func run(username string, fetcher port.CommitFetcher, debug bool) {
	fmt.Println("ℹ️  Using public GitHub API")

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
