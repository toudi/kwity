package common

import "time"

type Rate struct {
	UnitPrice
	StartDate string `yaml:"start-date" json:"start-date"`
	StartTime time.Time
	EndTime   time.Time
}
