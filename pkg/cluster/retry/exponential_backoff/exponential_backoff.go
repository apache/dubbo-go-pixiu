package exponential_backoff

import (
	"fmt"
	"github.com/apache/dubbo-go-pixiu/pkg/cluster/retry"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/mitchellh/mapstructure"
	"math"
	"math/rand"
	"time"
)

func init() {
	retry.RegisterRetryPolicy(model.RetryerExponentialBackoff, newExponentialBackoffRetry)
}

type ExponentialBackoffRetry struct {
	MaxAttempts     uint
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
	currentTry      uint
}

type ExponentialBackoffConfig struct {
	Times           uint    `mapstructure:"times" default:"3"`
	InitialInterval string  `mapstructure:"initialInterval" default:"100ms"`
	MaxInterval     string  `mapstructure:"maxInterval" default:"5s"`
	Multiplier      float64 `mapstructure:"multiplier" default:"2.0"`
}

func (e *ExponentialBackoffRetry) Attempt(err error) bool {
	if e.currentTry >= e.MaxAttempts {
		return false
	}

	// Don't wait before the first try
	if e.currentTry > 0 {
		backoff := float64(e.InitialInterval) * math.Pow(e.Multiplier, float64(e.currentTry-1))
		cappedBackoff := time.Duration(math.Min(backoff, float64(e.MaxInterval)))
		// Add jitter to prevent thundering herd
		jitter := time.Duration(rand.Intn(100)) * time.Millisecond
		time.Sleep(cappedBackoff + jitter)
	}

	e.currentTry++
	return true
}

func (e *ExponentialBackoffRetry) Reset() {
	e.currentTry = 0
}

func newExponentialBackoffRetry(config map[string]any) (retry.Retryer, error) {
	var cfg ExponentialBackoffConfig
	if err := mapstructure.Decode(config, &cfg); err != nil {
		return nil, fmt.Errorf("failed to decode exponential backoff config: %w", err)
	}

	initial, err := time.ParseDuration(cfg.InitialInterval)
	if err != nil {
		return nil, fmt.Errorf("invalid initialInterval: %w", err)
	}
	duration, err := time.ParseDuration(cfg.MaxInterval)
	if err != nil {
		return nil, fmt.Errorf("invalid maxInterval: %w", err)
	}

	return &ExponentialBackoffRetry{
		MaxAttempts:     cfg.Times + 1,
		InitialInterval: initial,
		MaxInterval:     duration,
		Multiplier:      cfg.Multiplier,
	}, nil
}
