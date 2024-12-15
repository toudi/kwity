package fs

import (
	"os"
	"path/filepath"

	"github.com/toudi/kwity/internal/db"
	"gopkg.in/yaml.v3"
)

var holidaysDBInstance *HolidaysDB

type HolidaysDB struct {
	BankHolidays []string `yaml:"bank-holidays"`
	PtoRanges    []string `yaml:"pto"`
}

func (b *HolidaysDB) GetBankHolidays() []string {
	return b.BankHolidays
}

func (b *HolidaysDB) GetPTORanges() []string {
	return b.PtoRanges
}

func (f *FSDb) Holidays() db.HolidaysInterface {
	if holidaysDBInstance == nil {
		holidaysDBInstance = &HolidaysDB{}

		holidaysFile, err := os.Open(filepath.Join(f.config.Root, "holidays.yaml"))
		if err == nil {
			_ = yaml.NewDecoder(holidaysFile).Decode(&holidaysDBInstance)
		}
	}

	return holidaysDBInstance
}
