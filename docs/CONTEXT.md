# DevLevel - Context

## Execution Environment

* The application runs locally on the user's machine
* Executed via terminal (CLI)
* No container required

---

## Language & Runtime

* Language: Go 1.21+
* Daily use: `devlevel` (compiled binary)
* From source: `go run ./cmd`

---

## External Integrations

* GitHub REST API v3 (public endpoints — no authentication required)
* Endpoints used:
  * `GET /users/{username}/events/public` — discovers repos with recent activity
  * `GET /repos/{owner}/{repo}/commits?author={username}&since={date}` — fetches commits per repo

---

## User Configuration

* GitHub username is saved locally via `devlevel setup`
* Config file location: `~/.devlevel/config.json`
* No environment variables required
* No tokens or authentication

---

## Local Persistence

* Progression state is saved at `~/.devlevel/state.json`
* Stores:
  * Total accumulated XP (never decreases)
  * Commit SHAs already awarded XP (prevents double-counting)
  * Days with activity (enables streak calculation beyond the API window)

---

## Constraints

* Use only Go's standard library (zero external dependencies)
* Keep code simple and clear
* Hexagonal architecture: domain isolated from infrastructure and presentation

---

## Known Limitations

* Only public repositories are counted (public API limitation)
* Public API rate limit: 60 requests/hour per IP — heavy use may hit the limit
* GitHub events API covers approximately the last 30 days of activity
