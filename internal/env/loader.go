package env

import (
	"bufio"
	"os"
	"strings"
)

// Load reads a .env file and sets any key=value pairs as environment variables.
// Variables already set in the environment take precedence (are not overwritten).
// Lines starting with # and blank lines are ignored.
func Load(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // .env is optional
		}
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		// Strip surrounding quotes: "value" or 'value' → value
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		// Environment variable already set takes precedence
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	return scanner.Err()
}
