package interfaces

import "github.com/toudi/kwity/internal/workdays"

type App interface {
	WorkdaysCalculator() (*workdays.WorkingDaysCalculator, error)
}
