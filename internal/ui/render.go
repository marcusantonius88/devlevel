package ui

import (
	"fmt"
	"strings"

	"devlevel/internal/model"
)

const (
	barWidth   = 10
	blockFull  = "█"
	blockEmpty = "░"
	labelWidth = 7 // pad label column so colons align
)

// Render prints the DevLevel stats to stdout.
//
// Layout:
//
//	🚀 DevLevel
//
//	🔥 CURRENT STREAK: 3 DAYS
//	✅ Daily Goal: COMPLETE
//
//	👤 User    : marcusantonius88
//	🏆 Level   : 1
//	⚡ XP      : 80
//	🏅 Rank    : Rookie
//
//	📈 Progress to Level 2
//	   ████████░░ 80%
//	🎯 Next Level: 20 XP remaining
//
//	📊 Summary
//	   • Recent activity: 8 commits
//	   • Keep the momentum going
func Render(s model.Stats, rank string, progressPct, xpRemaining int) {
	fmt.Println("🚀 DevLevel")

	// --- Streak block (hero element) ---
	fmt.Println()
	fmt.Printf("🔥 CURRENT STREAK: %s\n", streakLabel(s.Streak))
	fmt.Println(dailyGoalLine(s.DailyGoalMet))

	// --- Identity block ---
	fmt.Println()
	field("👤", "User", s.Username)
	field("🏆", "Level", fmt.Sprintf("%d", s.Level))
	field("⚡", "XP", fmt.Sprintf("%d", s.XP))
	field("🏅", "Rank", rank)

	// --- Progress block ---
	fmt.Println()
	if xpRemaining > 0 {
		fmt.Printf("📈 Progress to Level %d\n", s.Level+1)
		fmt.Printf("   %s %d%%\n", progressBar(progressPct), progressPct)
		fmt.Printf("🎯 Next Level: %d XP remaining\n", xpRemaining)
	} else {
		fmt.Println("📈 Progress: MAX LEVEL")
	}

	// --- Summary block ---
	fmt.Println()
	fmt.Println("📊 Summary")
	fmt.Printf("   • Recent activity: %s\n", formatCommits(s.CommitCount))
	fmt.Printf("   • %s\n", motivationalMessage(s.DailyGoalMet, s.Streak))
}

// streakLabel formats the streak count as "N DAY" / "N DAYS" in uppercase.
func streakLabel(n int) string {
	if n == 1 {
		return "1 DAY"
	}
	return fmt.Sprintf("%d DAYS", n)
}

// dailyGoalLine returns the daily goal status line.
func dailyGoalLine(met bool) string {
	if met {
		return "✅ Daily Goal: COMPLETE"
	}
	return "⚠️  Daily Goal: PENDING — commit today to protect your streak"
}

// motivationalMessage returns a context-aware message for the summary block.
func motivationalMessage(dailyGoalMet bool, streak int) string {
	if !dailyGoalMet {
		return "Commit today to protect your streak"
	}
	if streak >= 30 {
		return "Incredible consistency — keep it up"
	}
	if streak >= 7 {
		return "Keep the momentum going"
	}
	return "Daily goal completed — see you tomorrow"
}

// field prints a labelled row with consistent emoji + padded label + colon.
func field(emoji, label, value string) {
	fmt.Printf("%s %-*s : %s\n", emoji, labelWidth, label, value)
}

// progressBar builds a textual bar like ████████░░ for a given percentage.
func progressBar(pct int) string {
	filled := pct * barWidth / 100
	if filled > barWidth {
		filled = barWidth
	}
	return strings.Repeat(blockFull, filled) + strings.Repeat(blockEmpty, barWidth-filled)
}

// formatCommits returns "1 commit" or "N commits" correctly pluralised.
func formatCommits(n int) string {
	if n == 1 {
		return "1 commit"
	}
	return fmt.Sprintf("%d commits", n)
}
