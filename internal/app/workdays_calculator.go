package app

import (
	"github.com/samber/lo"
	"github.com/toudi/kwity/internal/workdays"
)

func (a *App) WorkdaysCalculator() (*workdays.WorkingDaysCalculator, error) {
	return workdays.NewWorkingDaysCalculator(
		lo.SliceToMap(a.db.Holidays().GetBankHolidays(), func(item string) (string, bool) {
			return item, true
		}),
		a.db.Holidays().GetPTORanges(),
	)
}
