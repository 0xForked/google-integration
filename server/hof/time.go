package hof

import (
	"fmt"
	"time"
)

// TimeToInt Get the current time for sql data
//
// usage:
//
//	currentTime := time.Now()
//	timeInt := TimeToInt(currentTime)
func TimeToInt(t time.Time) int {
	formattedTime := t.Format("15:04") // t.Format("03:04pm")
	var hour, minute int
	_, _ = fmt.Sscanf(formattedTime, "%d:%d", &hour, &minute)
	timeInt := hour*100 + minute
	return timeInt
}
