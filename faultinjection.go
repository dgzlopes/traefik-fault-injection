// Package traefik_fault_injection can inject faults via HTTP headers
package traefik_fault_injection

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// Config the plugin configuration.
type Config struct {
	Delay           bool `yaml:"Delay"`
	DelayDuration   int  `yaml:"DelayDuration"`
	DelayPercentage int  `yaml:"DelayPercentage"`
	Abort           bool `yaml:"Abort"`
	AbortCode       int  `yaml:"AbortCode"`
	AbortPercentage int  `yaml:"AbortPercentage"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Delay:           true,
		DelayDuration:   0,
		DelayPercentage: 100,
		Abort:           true,
		AbortCode:       400,
		AbortPercentage: 100,
	}
}

// FaultInjection plugin
type FaultInjection struct {
	next            http.Handler
	Delay           bool
	DelayDuration   int
	DelayPercentage int
	Abort           bool
	AbortCode       int
	AbortPercentage int
	name            string
}

// New created a new plugin
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &FaultInjection{
		Delay:           config.Delay,
		DelayDuration:   config.DelayDuration,
		DelayPercentage: config.DelayPercentage,
		Abort:           config.Abort,
		AbortCode:       config.AbortCode,
		AbortPercentage: config.DelayPercentage,
		next:            next,
		name:            name,
	}, nil
}

func (a *FaultInjection) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if a.Delay {
		delayHeader := req.Header.Get("X-Traefik-Fault-Delay-Request")
		delayPercentageHeader := req.Header.Get("X-Traefik-Fault-Delay-Request-Percentage")
		if FaultShouldRun(ParseHeaderValue(delayPercentageHeader, a.DelayPercentage)) {
			time.Sleep(time.Duration(ParseHeaderValue(delayHeader, a.DelayDuration)) * time.Millisecond)
		}
	}

	if a.Abort {
		abortHeader := req.Header.Get("X-Traefik-Fault-Abort-Request")
		abortPercentageHeader := req.Header.Get("X-Traefik-Fault-Abort-Request-Percentage")
		if len(abortHeader) != 0 {
			if FaultShouldRun(ParseHeaderValue(abortPercentageHeader, a.AbortPercentage)) {
				rw.WriteHeader(ParseHeaderValue(abortHeader, a.AbortCode))
				return
			}
		}
	}
	a.next.ServeHTTP(rw, req)
}

// ParseHeaderValue is used for the value transformation
func ParseHeaderValue(rawValue string, defaultValue int) int {
	if len(rawValue) != 0 {
		if parsedValue, err := strconv.Atoi(rawValue); err == nil {
			return parsedValue
		}
	}
	return defaultValue
}

// FaultShouldRun is used to check if the fault should run or not
func FaultShouldRun(percent int) bool {
	return rand.Intn(100) <= percent
}
