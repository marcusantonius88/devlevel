package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"devlevel/internal/config"
	"devlevel/internal/gamification"
	githubadapter "devlevel/internal/github"
	"devlevel/internal/model"
	"devlevel/internal/port"
	"devlevel/internal/ui"
)

func main() {
	debug := flag.Bool("debug", false, "Print debug information about API calls")
	flag.Parse()

	// Subcommand routing: "devlevel setup" or "devlevel"
	if flag.NArg() > 0 && flag.Arg(0) == "setup" {
		runSetup()
		return
	}

	cfg, err := config.Load()
	if err != nil {
		if errors.Is(err, config.ErrNotConfigured) {
			fmt.Fprintln(os.Stderr, "❌ No GitHub username configured.")
			fmt.Fprintln(os.Stderr, "Please run:")
			fmt.Fprintln(os.Stderr, "   devlevel setup")
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	adapter := githubadapter.NewClient()
	run(cfg.GitHubUsername, adapter, *debug)
}

// runSetup prompts the user for their GitHub username and saves it locally.
func runSetup() {
	fmt.Print("Enter your GitHub username: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	username := strings.TrimSpace(scanner.Text())

	if err := config.ValidateUsername(username); err != nil {
		fmt.Fprintln(os.Stderr, "❌ Invalid username:", err)
		os.Exit(1)
	}

	cfg := &config.Config{GitHubUsername: username}
	if err := config.Save(cfg); err != nil {
		fmt.Fprintln(os.Stderr, "❌ Could not save configuration:", err)
		os.Exit(1)
	}

	dir, _ := config.Dir()
	fmt.Println("✅ Configuration saved successfully")
	fmt.Printf("   Config location: %s/config.json\n", dir)
	fmt.Println()
	fmt.Println("You're all set. Run devlevel to check your streak.")
}

// run contains the application flow and depends only on port interfaces,
// making it straightforward to test with mock adapters.
func run(username string, fetcher port.CommitFetcher, debug bool) {
	fmt.Println("ℹ️  Using public GitHub API")

	result, err := fetcher.FetchRecentCommits(username, debug)
	if err != nil {
		if errors.Is(err, port.ErrRateLimit) {
			fmt.Fprintln(os.Stderr, "")
			fmt.Fprintln(os.Stderr, "⚠️  GitHub API rate limit reached.")
			fmt.Fprintln(os.Stderr, "   The public API allows 60 requests/hour without authentication.")
			fmt.Fprintln(os.Stderr, "   Please wait a few minutes and try again.")
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	// All repos timed out — data is entirely unreliable, don't show stats.
	if result.TotalRepos > 0 && result.SkippedRepos == result.TotalRepos {
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "⚠️  Could not retrieve your activity data.")
		fmt.Fprintln(os.Stderr, "   GitHub API is responding slowly right now.")
		fmt.Fprintln(os.Stderr, "   Your streak is safe — please try again in a few minutes.")
		os.Exit(1)
	}

	stats := buildStats(username, result.Commits)

	// Some repos timed out — warn that results may be incomplete.
	if result.SkippedRepos > 0 {
		fmt.Printf("⚠️  %d of %d repo(s) could not be reached — data may be incomplete\n\n",
			result.SkippedRepos, result.TotalRepos)
	}

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
