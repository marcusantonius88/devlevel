package model

import "time"

// Commit represents a single Git commit with its push timestamp.
type Commit struct {
	SHA  string
	Date time.Time
}

// Stats holds all computed metrics for a user.
type Stats struct {
	Username     string
	Level        int
	XP           int
	Streak       int
	CommitCount  int
	DailyGoalMet bool // true if the user has at least one commit today
}

// ----------------------------------------------------------------------------
// Extension points — not yet implemented, reserved for future features.
// ----------------------------------------------------------------------------

// Achievement represents a badge or milestone earned by the user.
// Reserved for a future achievements system.
type Achievement struct {
	ID          string
	Title       string
	Description string
	EarnedAt    time.Time
}

// Milestone marks a streak length worth celebrating (e.g. 7, 30, 100 days).
// Reserved for a future milestones system.
type Milestone struct {
	Days    int
	Label   string
	Reached bool
}
