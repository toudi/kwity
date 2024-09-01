package workdays_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/toudi/kwity/internal/workdays"
)

func TestWorkdays(t *testing.T) {
	type testCase struct {
		label                string
		inputStartDate       time.Time
		inputEndDate         time.Time
		inputBankHolidays    map[string]bool
		expectedWorkdays     int
		expectedBankHolidays int
	}

	for _, test := range []testCase{
		{
			label:                "Jan 1 -> Jan 2",
			inputStartDate:       time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
			inputEndDate:         time.Date(2024, 1, 2, 0, 0, 0, 0, time.Local),
			expectedWorkdays:     2,
			expectedBankHolidays: 0,
		},
		{
			label:                "Jan 1 -> Jan 8",
			inputStartDate:       time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
			inputEndDate:         time.Date(2024, 1, 8, 0, 0, 0, 0, time.Local),
			expectedWorkdays:     6,
			expectedBankHolidays: 0,
		},
		{
			label:                "Jan 1 -> Jan 31",
			inputStartDate:       time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
			inputEndDate:         time.Date(2024, 1, 31, 0, 0, 0, 0, time.Local),
			expectedWorkdays:     23,
			expectedBankHolidays: 0,
		},
		{
			label:          "May 1 -> May 31 with bank holidays",
			inputStartDate: time.Date(2024, 5, 1, 0, 0, 0, 0, time.Local),
			inputEndDate:   time.Date(2024, 5, 31, 0, 0, 0, 0, time.Local),
			inputBankHolidays: map[string]bool{
				"2024-05-01": true,
				"2024-05-03": true,
				"2024-05-30": true,
			},
			expectedWorkdays:     23, // there are 23 working days, however ..
			expectedBankHolidays: 3,  // .. within them there are 3 public holidays.
		},
	} {
		test := test
		t.Run(test.label, func(t *testing.T) {
			t.Parallel()

			numWorkdays, numBankHolidays := workdays.CalculateWorkingDays(
				test.inputStartDate, test.inputEndDate, test.inputBankHolidays,
			)

			assert.Equal(t, test.expectedWorkdays, numWorkdays)
			assert.Equal(t, test.expectedBankHolidays, numBankHolidays)
		})
	}
}
