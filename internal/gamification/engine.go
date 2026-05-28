package gamification

import (
	"time"

	"devlevel/internal/model"
)

// CalculateXP returns total XP: each commit is worth 10 points.
func CalculateXP(commits []model.Commit) int {
	return len(commits) * 10
}

// CalculateLevel maps XP to a level using the progression table from SPEC.md.
func CalculateLevel(xp int) int {
	switch {
	case xp >= 500:
		return 4
	case xp >= 250:
		return 3
	case xp >= 100:
		return 2
	default:
		return 1
	}
}

// LevelProgress returns the XP progress percentage towards the next level
// (0–100). Returns 100 if the user is already at max level.
func LevelProgress(xp int) int {
	type threshold struct{ from, to int }
	levels := []threshold{
		{0, 100},   // level 1 → 2
		{100, 250}, // level 2 → 3
		{250, 500}, // level 3 → 4
	}
	for _, l := range levels {
		if xp < l.to {
			return (xp - l.from) * 100 / (l.to - l.from)
		}
	}
	return 100 // max level
}

// XPToNextLevel returns how many XP are still needed to reach the next level.
// Returns 0 if the user is already at max level.
func XPToNextLevel(xp int) int {
	switch {
	case xp < 100:
		return 100 - xp
	case xp < 250:
		return 250 - xp
	case xp < 500:
		return 500 - xp
	default:
		return 0
	}
}

// RankTitle returns the title associated with a given level.
func RankTitle(level int) string {
	switch level {
	case 1:
		return "Rookie"
	case 2:
		return "Builder"
	case 3:
		return "Engineer"
	case 4:
		return "Architect"
	default:
		return "Architect"
	}
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
