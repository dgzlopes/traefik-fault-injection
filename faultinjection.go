// Package traefik_fault_injection can inject faults via HTTP headers
package traefik_fault_injection

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

// Config the plugin configuration.
type Config struct {
	delay        bool
	defaultDelay int
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		delay:        true,
		defaultDelay: 0,
	}
}

// FaultInjection plugin
type FaultInjection struct {
	next         http.Handler
	delay        bool
	defaultDelay int
	name         string
}

// New created a new plugin
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &FaultInjection{
		delay:        config.delay,
		defaultDelay: config.defaultDelay,
		next:         next,
		name:         name,
	}, nil
}

func (a *FaultInjection) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if a.delay == true {
		delayHeader := req.Header.Get("X-Traefik-Fault-Delay-Request")
		time.Sleep(time.Duration(ParseHeaderValue(delayHeader, a.defaultDelay)) * time.Millisecond)
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
