package broker

import (
	"fmt"
	"math"
	"time"
)

// roundToInt rounds a number to the nearest integer.
func roundToInt[T float32 | float64](f T) int {
	return int(math.Round(float64(f)))
}

// calculateDaysCount calculates the number of days between pickup and dropoff dates, inclusive, it adds a day if the dropoff time is after the pickup time to account for partial days.
func calculateDaysCount(pickupDate, pickupTime, dropoffDate, dropoffTime string) (int, error) {
	pickupDateTime, err := time.Parse("2006-01-02 15:04", pickupDate+" "+pickupTime)
	if err != nil {
		return 0, fmt.Errorf("invalid pickup datetime: %s %s: %w", pickupDate, pickupTime, err)
	}

	dropoffDateTime, err := time.Parse("2006-01-02 15:04", dropoffDate+" "+dropoffTime)
	if err != nil {
		return 0, fmt.Errorf("invalid dropoff datetime: %s %s: %w", dropoffDate, dropoffTime, err)
	}

	days := int(dropoffDateTime.Sub(pickupDateTime).Hours()/24) + 1
	if dropoffDateTime.Sub(pickupDateTime).Hours()/24 > float64(days) {
		days++
	}

	return days, nil
}
