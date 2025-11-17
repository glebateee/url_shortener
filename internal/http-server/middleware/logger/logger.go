package logger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "middleware/logger"),
		)
		log.Info("logger middleware enabled")

		fn := func(writer http.ResponseWriter, request *http.Request) {
			entry := log.With(
				slog.String("method", request.Method),
				slog.String("path", request.URL.Path),
				slog.String("remote_addr", request.RemoteAddr),
				slog.String("user_agent", request.UserAgent()),
				slog.String("request_id", middleware.GetReqID(request.Context())),
			)
			ww := middleware.NewWrapResponseWriter(writer, request.ProtoMajor)

			t1 := time.Now()

			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()
			next.ServeHTTP(ww, request)
		}
		return http.HandlerFunc(fn)
	}
}
