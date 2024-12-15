package workdays_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestWorkingDaysCalculator(t *testing.T) {
	type testCase struct {
		label              string
		inputStartDate     time.Time
		inputEndDate       time.Time
		inputBankHolidays  map[string]bool
		inputPtoRanges     []string
		inputLastDayCounts bool
		expected           workdays.WorkingDaysResult
		expectedError      error
	}

	for _, test := range []testCase{
		{
			label:              "no bank holidays, last day does not count",
			inputStartDate:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
			inputEndDate:       time.Date(2024, 1, 2, 0, 0, 0, 0, time.Local),
			inputLastDayCounts: false,
			expected: workdays.WorkingDaysResult{
				WorkingDays:    1,
				LastWorkingDay: time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
			},
		},
		{
			label:              "no bank holidays, last day counts",
			inputStartDate:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
			inputEndDate:       time.Date(2024, 1, 2, 0, 0, 0, 0, time.Local),
			inputLastDayCounts: true,
			expected: workdays.WorkingDaysResult{
				WorkingDays:    2,
				LastWorkingDay: time.Date(2024, 1, 2, 0, 0, 0, 0, time.Local),
			},
		},
		{
			label:          "no bank holidays, 1 PTO day defined as a single string, last day does not count",
			inputStartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
			inputEndDate:   time.Date(2024, 1, 2, 0, 0, 0, 0, time.Local),
			inputPtoRanges: []string{
				"2024-01-01",
			},
			expected: workdays.WorkingDaysResult{
				WorkingDays:    1,
				PtoDays:        1,
				LastWorkingDay: time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
			},
		},
		{
			label:              "no bank holidays, 1 PTO day defined as a single string, last day does count",
			inputStartDate:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.Local),
			inputEndDate:       time.Date(2024, 1, 2, 0, 0, 0, 0, time.Local),
			inputLastDayCounts: true,
			inputPtoRanges: []string{
				"2024-01-01",
			},
			expected: workdays.WorkingDaysResult{
				WorkingDays:    2,
				PtoDays:        1,
				LastWorkingDay: time.Date(2024, 1, 2, 0, 0, 0, 0, time.Local),
			},
		},
		{
			label:          "no bank holidays, single weekend",
			inputStartDate: time.Date(2024, 11, 04, 0, 0, 0, 0, time.Local),
			inputEndDate:   time.Date(2024, 11, 10, 0, 0, 0, 0, time.Local),
			expected: workdays.WorkingDaysResult{
				WorkingDays:    5,
				LastWorkingDay: time.Date(2024, 11, 8, 0, 0, 0, 0, time.Local),
			},
		},
		{
			label:             "single bank holiday, single weekend",
			inputStartDate:    time.Date(2024, 11, 11, 0, 0, 0, 0, time.Local),
			inputEndDate:      time.Date(2024, 11, 17, 0, 0, 0, 0, time.Local),
			inputBankHolidays: map[string]bool{"2024-11-11": true},
			expected: workdays.WorkingDaysResult{
				WorkingDays:    4,
				BankHolidays:   1,
				LastWorkingDay: time.Date(2024, 11, 15, 0, 0, 0, 0, time.Local),
			},
		},
	} {
		t.Run(test.label, func(t *testing.T) {
			test := test
			calculator, err := workdays.NewWorkingDaysCalculator(test.inputBankHolidays, test.inputPtoRanges)
			require.Equal(t, test.expectedError, err)
			if err == nil {
				result := calculator.CalculateWorkingDays(test.inputStartDate, test.inputEndDate, test.inputLastDayCounts)
				require.Equal(t, test.expected, result)
			}
		})
	}
}
