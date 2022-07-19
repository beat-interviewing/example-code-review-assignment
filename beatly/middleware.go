package beatly

import (
	"net/http"
	"time"

	"github.com/go-kit/log"
)

func LoggingMiddleware(logger log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			defer func(start time.Time) {
				next.ServeHTTP(w, req)
				now := time.Now()
				logger.Log(
					"ts", now,
					"d", now.Sub(start),
					"method", req.Method,
					"path", req.URL.Path,
					"ip", req.RemoteAddr,
				)
			}(time.Now())
		})
	}
}
