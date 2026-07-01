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
		fmt.Println()
	} else {
		fmt.Println("📈 Progress: MAX LEVEL")
		fmt.Println()
	}
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
