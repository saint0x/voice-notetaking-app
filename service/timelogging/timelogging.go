// service/timelogging/timelogging.go

package timelogging

import (
  "errors"
  "fmt"
  "time"
)

var startTime time.Time

// StartLogging starts the time logging for the given action.
func StartLogging() {
  startTime = time.Now()
  fmt.Printf("Started logging at: %s\n", startTime.Format("2006-01-02 15:04:05"))
}

// EndLogging logs the end time and calculates the duration for the given action.
func EndLogging(action string) {
  endTime := time.Now()
  duration := endTime.Sub(startTime)
  fmt.Printf("Finished %s at: %s\n", action, endTime.Format("2006-01-02 15:04:05"))
  fmt.Printf("Duration for %s: %s\n", action, duration)
}

// LogTime logs time for given action.
func LogTime(action string) error {
  if startTime.IsZero() {
    return errors.New("time logging not started")
  }

  EndLogging(action)
  return nil
}
