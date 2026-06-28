# DevLevel - Specification

## Overview

DevLevel is a CLI tool that helps developers build consistency through daily GitHub commit streaks. Inspired by Duolingo's streak mechanic, it turns the habit of committing into a gamified experience with streak tracking, daily goals, XP progression, and rank titles.

---

## Feature: Initial Setup

GIVEN the user runs `devlevel setup` for the first time
WHEN the system prompts for a GitHub username
THEN it should save the username to `~/.devlevel/config.json`
AND confirm with "✅ Configuration saved successfully"

GIVEN the user runs `devlevel` without having done setup
WHEN the system finds no local configuration
THEN it should display "❌ No GitHub username configured."
AND guide the user to run `devlevel setup`

---

## Feature: Fetch Activity

GIVEN the username is configured
WHEN the system runs
THEN it should discover active repos via `/users/{username}/events/public`
AND for each repo, fetch commits by the author via `/repos/{owner}/{repo}/commits`
AND only consider commits from the last 30 days

GIVEN the API returns a 403 error
WHEN the system detects a rate limit
THEN it should display an explanatory message about the 60 requests/hour limit
AND guide the user to try again in a few minutes

GIVEN one or more repos do not respond within the timeout
WHEN all repos fail
THEN it should display "Your streak is safe — please try again in a few minutes"
AND not show zeroed stats that could confuse the user

GIVEN only some repos fail due to timeout
WHEN partial data is retrieved
THEN it should display an incomplete data warning
AND show stats with the available data

---

## Feature: XP Calculation

GIVEN a list of new commits (SHAs not yet awarded XP)
WHEN the system processes the commits
THEN each new commit should add 10 XP to the accumulated total
AND the total XP should be persisted in `~/.devlevel/state.json`
AND XP should never decrease between runs

Example:

* 5 new commits → +50 XP

---

## Feature: Level Calculation (Craft Track)

GIVEN a total accumulated XP value
WHEN the system calculates the level
THEN it should return the level corresponding to the Craft Track table:

| Level | Rank        | Min XP |
|-------|-------------|--------|
| 1     | Apprentice  | 0      |
| 2     | Craftsman   | 100    |
| 3     | Artisan     | 250    |
| 4     | Forger      | 500    |
| 5     | Blacksmith  | 750    |
| 6     | Grandmaster | 1000   |
| 7     | Sage        | 1500   |
| 8     | Oracle      | 2000   |
| 9     | Mythic      | 3000   |

---

## Feature: Streak Calculation

GIVEN the active days persisted in the local state
WHEN the system calculates the streak
THEN it should count consecutive days with at least one commit, going backwards from today
AND if there is no commit today, it should start counting from yesterday (mid-day streak protection)
AND it should use the local timezone of the user's machine

Example:

* Activity on 3 consecutive days → Streak = 3
* Missing one day → streak resets to zero

---

## Feature: Daily Goal

GIVEN the commits for the day have been processed
WHEN the system checks the daily goal
THEN it should display "✅ Daily Goal: COMPLETE" if there is at least one commit today (local timezone)
AND it should display "⚠️ Daily Goal: PENDING — commit today to protect your streak" otherwise

---

## Feature: CLI Output

GIVEN the data has been processed
WHEN the system displays the information
THEN the output should follow the layout below, with streak as the main element:

```
🚀 DevLevel
ℹ️  Using public GitHub API

🔥 CURRENT STREAK: N DAYS
✅ Daily Goal: COMPLETE

👤 User    : <username>
🏆 Level   : <level>
⚡ XP      : <xp>
🏅 Rank    : <rank>

📈 Progress to Level N
   ██████░░░░ N%
🎯 Next Level: N XP remaining

📊 Summary
   • Last 30 days: N commits
   • <contextual motivational message>
```

Motivational message adapts to context:
* Streak < 7 days → "Daily goal completed — see you tomorrow"
* Streak 7–29 days → "Keep the momentum going"
* Streak 30+ days → "Incredible consistency — keep it up"
* Daily goal pending → "Commit today to protect your streak"

---

## Feature: Error Handling

GIVEN the user has not run setup
WHEN the system starts
THEN it should display a message guiding the user to run `devlevel setup`

GIVEN the API returns a rate limit error (403)
WHEN the system detects it
THEN it should display a friendly message explaining the limit and asking to try again later

GIVEN all repos hit timeout
WHEN no data is retrieved
THEN it should inform the user that it is a communication issue and that the streak is safe

---

## Out of Scope (MVP)

* No external database
* No frontend
* No multi-user support
* No private repository support (public API limitation)
* No OAuth or token-based authentication
