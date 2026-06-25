package gamification

import (
	"time"

	"devlevel/internal/model"
)

// CalculateXP returns total XP: each commit is worth 10 points.
func CalculateXP(commits []model.Commit) int {
	return len(commits) * 10
}

// levels defines the full progression table for the Craft track.
var levels = []struct {
	minXP int
	rank  string
}{
	{0, "Apprentice"},
	{100, "Craftsman"},
	{250, "Artisan"},
	{500, "Forger"},
	{750, "Blacksmith"},
	{1000, "Grandmaster"},
	{1500, "Sage"},
	{2000, "Oracle"},
	{3000, "Mythic"},
}

// CalculateLevel returns the level (1-based) for the given XP.
func CalculateLevel(xp int) int {
	level := 1
	for i, l := range levels {
		if xp >= l.minXP {
			level = i + 1
		}
	}
	return level
}

// LevelProgress returns the XP progress percentage towards the next level
// (0–100). Returns 100 if the user is already at max level.
func LevelProgress(xp int) int {
	for i, l := range levels {
		if i == len(levels)-1 {
			return 100 // max level
		}
		next := levels[i+1]
		if xp < next.minXP {
			return (xp - l.minXP) * 100 / (next.minXP - l.minXP)
		}
	}
	return 100
}

// XPToNextLevel returns how many XP are still needed to reach the next level.
// Returns 0 if the user is already at max level.
func XPToNextLevel(xp int) int {
	for i, l := range levels {
		_ = l
		if i == len(levels)-1 {
			return 0 // max level
		}
		next := levels[i+1]
		if xp < next.minXP {
			return next.minXP - xp
		}
	}
	return 0
}

// RankTitle returns the Craft track rank title for a given level.
func RankTitle(level int) string {
	if level < 1 || level > len(levels) {
		return levels[len(levels)-1].rank
	}
	return levels[level-1].rank
}

// DailyGoalMet returns true if the user has at least one commit today,
// evaluated in the local timezone of the machine running the tool.
func DailyGoalMet(commits []model.Commit) bool {
	today := time.Now().Format("2006-01-02") // local time
	for _, c := range commits {
		if c.Date.Local().Format("2006-01-02") == today {
			return true
		}
	}
	return false
}

// CalculateStreak counts consecutive days with at least one commit going
// backwards from today, evaluated in the local timezone of the machine.
// If there's no activity today yet, it starts from yesterday to avoid
// breaking an active streak mid-day.
func CalculateStreak(commits []model.Commit) int {
	if len(commits) == 0 {
		return 0
	}

	// Map active days using local time so the user's timezone is respected.
	activeDays := make(map[string]bool, len(commits))
	for _, c := range commits {
		activeDays[c.Date.Local().Format("2006-01-02")] = true
	}

	today := time.Now() // local time
	start := today
	if !activeDays[today.Format("2006-01-02")] {
		// No commit yet today — start counting from yesterday
		start = today.AddDate(0, 0, -1)
	}

	streak := 0
	current := start
	for activeDays[current.Format("2006-01-02")] {
		streak++
		current = current.AddDate(0, 0, -1)
	}

	return streak
}

// ----------------------------------------------------------------------------
// Extension points — stubs reserved for future features.
// ----------------------------------------------------------------------------

// Milestones returns which streak milestones the user has reached.
// Not yet implemented — reserved for a future milestones system.
func Milestones(_ int) []model.Milestone {
	return nil
}

// Achievements returns badges earned by the user based on their stats.
// Not yet implemented — reserved for a future achievements system.
func Achievements(_ model.Stats) []model.Achievement {
	return nil
}
