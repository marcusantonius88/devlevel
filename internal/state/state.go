// Package state manages the persisted progression data for a DevLevel user.
// It lives at ~/.devlevel/state.json and tracks:
//   - Cumulative XP (by recording seen commit SHAs — never decreases)
//   - Active days (by recording dates with at least one commit — used for
//     streak calculation so streaks longer than the API window are preserved)
package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"devlevel/internal/model"
)

const (
	fileName    = "state.json"
	xpPerCommit = 10
)

// State holds the persisted progression for a user.
type State struct {
	TotalXP     int             `json:"total_xp"`
	SeenCommits map[string]bool `json:"seen_commits"` // SHA → already awarded XP
	ActiveDays  map[string]bool `json:"active_days"`  // "2006-01-02" → had commit
}

// New returns a zero-value State ready to use.
func New() *State {
	return &State{
		SeenCommits: make(map[string]bool),
		ActiveDays:  make(map[string]bool),
	}
}

// Load reads state from the given directory.
// Returns a fresh State if the file does not exist yet.
func Load(dir string) (*State, error) {
	path := filepath.Join(dir, fileName)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return New(), nil
		}
		return nil, fmt.Errorf("could not read state: %w", err)
	}

	s := New()
	if err := json.Unmarshal(data, s); err != nil {
		return nil, fmt.Errorf("could not parse state: %w", err)
	}

	// Guard against nil maps from a malformed file.
	if s.SeenCommits == nil {
		s.SeenCommits = make(map[string]bool)
	}
	if s.ActiveDays == nil {
		s.ActiveDays = make(map[string]bool)
	}

	return s, nil
}

// Save writes state to the given directory.
func Save(dir string, s *State) error {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("could not create state directory: %w", err)
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("could not encode state: %w", err)
	}

	path := filepath.Join(dir, fileName)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("could not write state: %w", err)
	}

	return nil
}

// ApplyCommits awards XP for commits not yet seen and records active days.
// Both operations update the state in place.
func (s *State) ApplyCommits(commits []model.Commit) {
	for _, c := range commits {
		// Record the active day (local time).
		day := c.Date.Local().Format("2006-01-02")
		s.ActiveDays[day] = true

		// Award XP only for new commits.
		if c.SHA == "" {
			// No SHA to deduplicate — award XP unconditionally.
			s.TotalXP += xpPerCommit
			continue
		}
		if !s.SeenCommits[c.SHA] {
			s.SeenCommits[c.SHA] = true
			s.TotalXP += xpPerCommit
		}
	}
}

// CalculateStreak counts consecutive days with activity going backwards from
// today, using the persisted active days map (not just the current API window).
// If there's no activity today yet, it starts from yesterday.
func (s *State) CalculateStreak() int {
	if len(s.ActiveDays) == 0 {
		return 0
	}

	today := time.Now()
	start := today
	if !s.ActiveDays[today.Format("2006-01-02")] {
		start = today.AddDate(0, 0, -1)
	}

	streak := 0
	current := start
	for s.ActiveDays[current.Format("2006-01-02")] {
		streak++
		current = current.AddDate(0, 0, -1)
	}

	return streak
}

// DailyGoalMet returns true if there is recorded activity for today.
func (s *State) DailyGoalMet() bool {
	return s.ActiveDays[time.Now().Format("2006-01-02")]
}

// ActiveDaysSorted returns all active days sorted descending, for display.
func (s *State) ActiveDaysSorted() []string {
	days := make([]string, 0, len(s.ActiveDays))
	for d := range s.ActiveDays {
		days = append(days, d)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(days)))
	return days
}
