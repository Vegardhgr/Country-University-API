package handler

import "time"

var startTime time.Time

// SetTime
/*Sets the start time of the server*/
func SetTime() {
	startTime = time.Now()
}

// GetTime
/*Gets the time since last service restart*/
func GetTime() float64 {
	upTime := time.Since(startTime).Seconds()
	return upTime
}
