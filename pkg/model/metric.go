package model

import (
	"time"
)

type Metric struct {
	Enable   bool          `yaml:"enable" json:"enable"`
	Interval time.Duration `yaml:"interval" json:"interval"`
}
