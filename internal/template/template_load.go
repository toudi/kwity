package template

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

var errContractorIsEmpty = errors.New("contractor not set")

func LoadFromFile(srcFile string) (*Template, error) {
	var template Template

	templateFile, err := os.Open(srcFile)
	if err != nil {
		return nil, fmt.Errorf("error opening template: %v", err)
	}
	defer templateFile.Close()

	// TODO: replace with template interface to support different formats.
	if strings.ToLower(path.Ext(srcFile)) == ".yaml" {
		decoder := yaml.NewDecoder(templateFile)
		if err = decoder.Decode(&template); err != nil {
			return nil, errors.Join(ErrTemplateParseError, err)
		}
	}

	// check if the contractor is not empty
	if template.Recipient == nil {
		return nil, errContractorIsEmpty
	}

	// parse rates
	if err = template.parseRatesTable(); err != nil {
		fmt.Printf("error parsing rates: %v\n", err)
		return nil, errInvalidRates
	}

	// // let's check if the quantity is ok
	// for i, item := range template.Items {
	// 	if !item.CalculateNumberOfWorkdays() {
	// 		// let's check if it's a valid number.
	// 		if err = item.parseQuantity(); err != nil {
	// 			return nil, errors.Join(
	// 				ErrInvalidQuantity,
	// 				fmt.Errorf("quantity at item %d is invalid: %v", i, item.Quantity),
	// 			)
	// 		}
	// 	}
	// }
	return &template, nil
}
