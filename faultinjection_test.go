package traefik_fault_injection

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServeHTTPp(t *testing.T) {
	tests := []struct {
		desc                    string
		delay                   bool
		delayDuration           int
		delayPercentage         int
		abort                   bool
		abortCode               int
		abortPercentage         int
		expDelayDurationSeconds float64
		expStatusCode           int
	}{
		{
			desc:                    "delay is disabled",
			delay:                   false,
			delayDuration:           5000,
			delayPercentage:         100,
			abort:                   false,
			abortCode:               400,
			abortPercentage:         0,
			expDelayDurationSeconds: 0.5,
			expStatusCode:           http.StatusOK,
		}, {
			desc:                    "delay is enabled",
			delay:                   true,
			delayDuration:           5000,
			delayPercentage:         100,
			abort:                   false,
			abortCode:               400,
			abortPercentage:         0,
			expDelayDurationSeconds: 5.5,
			expStatusCode:           http.StatusOK,
		}, {
			desc:                    "delay is enabled but percentage is 0",
			delay:                   true,
			delayDuration:           5000,
			delayPercentage:         0,
			abort:                   false,
			abortCode:               400,
			abortPercentage:         0,
			expDelayDurationSeconds: 0.5,
			expStatusCode:           http.StatusOK,
		}, {
			desc:                    "abort is enabled but percentage is 0",
			delay:                   false,
			delayDuration:           5000,
			delayPercentage:         0,
			abort:                   true,
			abortCode:               400,
			abortPercentage:         0,
			expDelayDurationSeconds: 0.5,
			expStatusCode:           http.StatusOK,
		}, {
			desc:                    "abort is enabled with code 400",
			delay:                   false,
			delayDuration:           5000,
			delayPercentage:         0,
			abort:                   true,
			abortCode:               400,
			abortPercentage:         100,
			expDelayDurationSeconds: 0.5,
			expStatusCode:           http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			cfg := &Config{
				Delay:           test.delay,
				DelayDuration:   test.delayDuration,
				DelayPercentage: test.delayPercentage,
				Abort:           test.abort,
				AbortCode:       test.abortCode,
				AbortPercentage: test.abortPercentage,
			}

			start := time.Now()

			next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

			handler, err := New(context.Background(), next, cfg, "fault")
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, "http://localhost", nil)

			handler.ServeHTTP(recorder, req)

			elapsed := time.Since(start)

			if elapsed.Seconds() > test.expDelayDurationSeconds {
				t.Errorf("got unexpected delay: enlapsed %f expected %f", elapsed.Seconds(), test.expDelayDurationSeconds)
			}

			if recorder.Result().StatusCode != test.expStatusCode {
				t.Errorf("got status code %d, want %d", recorder.Code, test.expStatusCode)
			}
		})
	}
}
