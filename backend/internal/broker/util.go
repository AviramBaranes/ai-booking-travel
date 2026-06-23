package broker

import (
	"fmt"
	"strings"
	"time"

	"encore.dev/rlog"
)

// CalculateDaysCount calculates the number of days between pickup and dropoff dates, inclusive, it adds a day if the dropoff time is after the pickup time to account for partial days.
func CalculateDaysCount(pickupDate, pickupTime, dropoffDate, dropoffTime string) (int, error) {
	pickupDateTime, err := time.Parse("2006-01-02 15:04", pickupDate+" "+pickupTime)
	if err != nil {
		return 0, fmt.Errorf("invalid pickup datetime: %s %s: %w", pickupDate, pickupTime, err)
	}

	dropoffDateTime, err := time.Parse("2006-01-02 15:04", dropoffDate+" "+dropoffTime)
	if err != nil {
		return 0, fmt.Errorf("invalid dropoff datetime: %s %s: %w", dropoffDate, dropoffTime, err)
	}

	days := int(dropoffDateTime.Sub(pickupDateTime).Hours() / 24)
	if dropoffDateTime.Sub(pickupDateTime).Hours()/24 > float64(days) {
		days++
	}

	return days, nil
}

// isElectric checks if the car is electric based on the ACRISS code
func isElectric(acrissCode string) bool {
	if len(acrissCode) < 4 {
		rlog.Warn("invalid acriss code, cannot determine if the car is electric", "acrissCode", acrissCode)
		return false
	}
	return acrissCode[3] == 'E' || acrissCode[3] == 'C'
}

// normalizeModelName removes the "or similar" suffix from the model name if it exists
func normalizeModelName(model string) string {
	const suffix = " or similar"

	if len(model) >= len(suffix) && strings.HasSuffix(strings.ToLower(model), suffix) {
		return strings.TrimSpace(model[:len(model)-len(suffix)])
	}

	return model
}

func isAcrissShowsAutoGear(acriss string) bool {
	if len(acriss) < 4 {
		rlog.Warn("invalid acriss code, cannot determine if the car is automatic", "acrissCode", acriss)
		return false
	}

	return acriss[2] == 'A' || acriss[2] == 'B' || acriss[2] == 'D'
}
