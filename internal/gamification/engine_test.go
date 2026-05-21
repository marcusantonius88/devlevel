package gamification_test

import (
	"testing"
	"time"

	"devlevel/internal/gamification"
	"devlevel/internal/model"
)

// commit builds a model.Commit with a date N days before today (UTC).
func commit(daysAgo int) model.Commit {
	return model.Commit{
		SHA:  "abc123",
		Date: time.Now().UTC().AddDate(0, 0, -daysAgo),
	}
}

// commits builds a slice of commits, one per day offset provided.
func commits(daysAgo ...int) []model.Commit {
	out := make([]model.Commit, len(daysAgo))
	for i, d := range daysAgo {
		out[i] = commit(d)
	}
	return out
}

// ── CalculateXP ──────────────────────────────────────────────────────────────

func TestCalculateXP_NoCommits(t *testing.T) {
	if got := gamification.CalculateXP(nil); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestCalculateXP_SingleCommit(t *testing.T) {
	if got := gamification.CalculateXP(commits(0)); got != 10 {
		t.Errorf("expected 10, got %d", got)
	}
}

func TestCalculateXP_MultipleCommits(t *testing.T) {
	// 5 commits → 50 XP (example from SPEC)
	if got := gamification.CalculateXP(commits(0, 1, 2, 3, 4)); got != 50 {
		t.Errorf("expected 50, got %d", got)
	}
}

// ── CalculateLevel ───────────────────────────────────────────────────────────

func TestCalculateLevel(t *testing.T) {
	tests := []struct {
		xp    int
		want  int
		label string
	}{
		{0, 1, "zero XP"},
		{99, 1, "just below level 2 threshold"},
		{100, 2, "exactly level 2 threshold"},
		{249, 2, "just below level 3 threshold"},
		{250, 3, "exactly level 3 threshold"},
		{499, 3, "just below level 4 threshold"},
		{500, 4, "exactly level 4 threshold"},
		{9999, 4, "very high XP stays at level 4"},
	}

	for _, tt := range tests {
		t.Run(tt.label, func(t *testing.T) {
			got := gamification.CalculateLevel(tt.xp)
			if got != tt.want {
				t.Errorf("CalculateLevel(%d) = %d, want %d", tt.xp, got, tt.want)
			}
		})
	}
}

// ── CalculateStreak ──────────────────────────────────────────────────────────

func TestCalculateStreak_NoCommits(t *testing.T) {
	if got := gamification.CalculateStreak(nil); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestCalculateStreak_OnlyToday(t *testing.T) {
	// Activity today → streak of 1
	if got := gamification.CalculateStreak(commits(0)); got != 1 {
		t.Errorf("expected 1, got %d", got)
	}
}

func TestCalculateStreak_OnlyYesterday(t *testing.T) {
	// No commit today but commit yesterday → streak still 1
	// (streak must not break just because the day hasn't ended yet)
	if got := gamification.CalculateStreak(commits(1)); got != 1 {
		t.Errorf("expected 1, got %d", got)
	}
}

func TestCalculateStreak_ConsecutiveDays(t *testing.T) {
	// today + 2 previous days → streak of 3
	if got := gamification.CalculateStreak(commits(0, 1, 2)); got != 3 {
		t.Errorf("expected 3, got %d", got)
	}
}

func TestCalculateStreak_GapBreaksStreak(t *testing.T) {
	// today and 2 days ago, but gap on day 1 → streak is 1 (only today)
	if got := gamification.CalculateStreak(commits(0, 2)); got != 1 {
		t.Errorf("expected 1, got %d", got)
	}
}

func TestCalculateStreak_OldCommitsOnlyNoStreak(t *testing.T) {
	// Last activity was 3 days ago with a gap before today → streak 0
	// days 3 and 5 ago: no consecutive chain reaching today/yesterday
	if got := gamification.CalculateStreak(commits(3, 5)); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestCalculateStreak_MultipleCommitsSameDay(t *testing.T) {
	// Two commits on the same day count as a single active day
	c := []model.Commit{commit(0), commit(0), commit(1)}
	if got := gamification.CalculateStreak(c); got != 2 {
		t.Errorf("expected 2, got %d", got)
	}
}
