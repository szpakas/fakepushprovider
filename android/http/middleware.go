package http

import (
	"net/http"
	"net/http/httptest"

	"github.com/uber-go/zap"
)

type LoggingMiddleware struct {
	Handler http.Handler
	Logger zap.Logger
}

func (m *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rec := httptest.NewRecorder()

	m.Handler.ServeHTTP(rec, r)

	m.Logger.Debug(
		"request:responded",
		zap.Int("res:status", rec.Code),
		zap.String("res:body", rec.Body.String()),
	)

	// -- recreate the response
	for k, v := range rec.Header() {
		w.Header()[k] = v
	}
	w.WriteHeader(rec.Code)
	w.Write(rec.Body.Bytes())
}
