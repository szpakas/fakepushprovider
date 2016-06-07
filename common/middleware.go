package common

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/uber-go/zap"
)

// LoggingMiddleware allows logging of the request and response.
type LoggingMiddleware struct {
	// Handler is the handler to be wrapped
	Handler http.Handler

	// Logger is the instance of zap.Logger used in logging
	Logger zap.Logger

	// reqCnt is an internal counter tracking total number of requests from the process launch.
	reqCnt uint64
}

// ServeHTTP call next middleware/handler but substitutes response writer with recorder to catch the response.
// Event is logged.
//
// Header X-requests-total is added with total number of requests handled from the process launch.
func (m *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rec := httptest.NewRecorder()

	m.Handler.ServeHTTP(rec, r)

	rT := atomic.AddUint64(&m.reqCnt, 1)
	rec.Header().Add("X-requests-total", fmt.Sprintf("%d", rT))

	// -- recreate the response
	for k, v := range rec.Header() {
		w.Header()[k] = v
	}
	w.WriteHeader(rec.Code)
	w.Write(rec.Body.Bytes())

	fields := []zap.Field{
		zap.Int("res:status", rec.Code),
		zap.String("res:body", rec.Body.String()),
	}
	for k, v := range r.Header {
		fields = append(fields, zap.String("req:header:"+k, strings.Join(v, ",")))
	}
	for k, v := range rec.Header() {
		fields = append(fields, zap.String("res:header:"+k, strings.Join(v, ",")))
	}
	m.Logger.Debug(
		"request:responded",
		fields...,
	)
}

// DelayMiddleware allows to delay response.
// May be used as poor man's network delay simulator.
type DelayMiddleware struct {
	// Handler is the handler to be wrapped
	Handler http.Handler

	// DelayFn is responsible for introducing the delay.
	// It's called before handler/next middleware.
	// Number of milliseconds delayed should be returned.
	//
	//  func() int {
	//      d := 20 + rand.Intn(100)
	//      time.Sleep(time.Millisecond * time.Duration(d))
	//      return d
	//  }
	DelayFn func() (delay int)
}

// ServeHTTP is introducing delay before next step in chain.
//
// X-delayed-by-ms header is added with number of milliseconds delayed.
func (m *DelayMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d := m.DelayFn()

	m.Handler.ServeHTTP(w, r)

	w.Header().Add("X-delayed-by-ms", strconv.Itoa(d))
}
