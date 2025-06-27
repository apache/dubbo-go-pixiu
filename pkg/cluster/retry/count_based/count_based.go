package count_based

import (
	"fmt"
	"github.com/apache/dubbo-go-pixiu/pkg/cluster/retry"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func init() {
	retry.RegisterRetryPolicy(model.RetryerCountBased, newCountBasedRetry)
}

type CountBasedRetry struct {
	MaxAttempts uint
	currentTry  uint
}

func (r *CountBasedRetry) Attempt(err error) bool {
	if r.currentTry < r.MaxAttempts {
		r.currentTry++
		return true
	}
	return false
}

func (r *CountBasedRetry) Reset() {
	r.currentTry = 0
}

func newCountBasedRetry(config map[string]any) (retry.Retryer, error) {
	timesValue, exists := config["times"]
	if !exists {
		return nil, fmt.Errorf("'times' field is missing in retry configuration")
	}

	timesUint, ok := timesValue.(int)
	if !ok {
		return nil, fmt.Errorf("invalid type for 'retry.count_based.times', expected int but got %T", timesValue)
	}

	// Total attempts = 1 initial try plus number of retries.
	return &CountBasedRetry{MaxAttempts: uint(timesUint) + 1}, nil
}
