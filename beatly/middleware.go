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
				iw := &interceptingWriter{w, http.StatusOK}
				next.ServeHTTP(iw, req)
				now := time.Now()
				logger.Log(
					"time", now,
					"duration", now.Sub(start),
					"method", req.Method,
					"path", req.URL.Path,
					"ip", req.RemoteAddr,
					"status", iw.code,
				)
			}(time.Now())
		})
	}
}

type interceptingWriter struct {
	http.ResponseWriter
	code int
}

// WriteHeader may not be explicitly called, so care must be taken to
// initialize w.code to its default value of http.StatusOK.
func (w *interceptingWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}
