package fs

import (
	"os"
	"path/filepath"

	"github.com/toudi/kwity/internal/db"
	"gopkg.in/yaml.v3"
)

var bankHolidaysDBInstance *BankHolidaysDB

type BankHolidaysDB struct {
	Holidays []string `yaml:"bank-holidays"`
}

func (b *BankHolidaysDB) GetBankHolidays() []string {
	return b.Holidays
}

func (f *FSDb) BankHolidays() db.BankHolidaysInterface {
	if bankHolidaysDBInstance == nil {
		bankHolidaysDBInstance = &BankHolidaysDB{}

		bankHolidaysFile, err := os.Open(filepath.Join(f.config.Root, "bank-holidays.yaml"))
		if err == nil {
			_ = yaml.NewDecoder(bankHolidaysFile).Decode(&bankHolidaysDBInstance)
		}
	}

	return bankHolidaysDBInstance
}
