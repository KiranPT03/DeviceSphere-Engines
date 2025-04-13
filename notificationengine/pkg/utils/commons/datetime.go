package commons

import (
	"time"
)

// GetCurrentUTCTimestamp returns the current UTC time as a string in RFC3339 format
// Example output: "2023-04-25T14:30:45Z"
func GetCurrentUTCTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// GetUTCTimestampFromTime converts a given time.Time to a UTC timestamp string
func GetUTCTimestampFromTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// GetUTCTimestampWithFormat returns the current UTC time as a string with the specified format
// Example: format = "2006-01-02 15:04:05" will return "2023-04-25 14:30:45"
func GetUTCTimestampWithFormat(format string) string {
	return time.Now().UTC().Format(format)
}

// ParseUTCTimestamp parses a timestamp string in RFC3339 format to time.Time
// Returns the parsed time and an error if parsing fails
func ParseUTCTimestamp(timestamp string) (time.Time, error) {
	return time.Parse(time.RFC3339, timestamp)
}

// AddDuration adds the specified duration to the current UTC time and returns the result as a string
// Example: AddDuration(24 * time.Hour) adds 24 hours to the current time
func AddDuration(duration time.Duration) string {
	return time.Now().UTC().Add(duration).Format(time.RFC3339)
}

// FormatTimestamp formats a timestamp string from one format to another
// Returns the formatted timestamp and an error if parsing fails
func FormatTimestamp(timestamp, inputFormat, outputFormat string) (string, error) {
	t, err := time.Parse(inputFormat, timestamp)
	if err != nil {
		return "", err
	}
	return t.UTC().Format(outputFormat), nil
}