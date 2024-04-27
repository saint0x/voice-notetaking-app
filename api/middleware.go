// internal/api/middleware.go

package api

import (
  "log"
  "net/http"
  "time"
)

// LoggerMiddleware logs incoming requests.
func LoggerMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    startTime := time.Now()

    next.ServeHTTP(w, r)

    log.Printf(
      "HTTP %s %s %s",
      r.Method,
      r.RequestURI,
      time.Since(startTime),
    )
  })
}
