package invoice

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/flosch/pongo2"
	"github.com/toudi/kwity/config"
)

const orderNumbersFile = "numbers.json"
const globalKey = "global"

func parseOrderNumbers() (map[string]interface{}, error) {
	var orderNumbers map[string]interface{} = make(map[string]interface{})
	orderNumbersFileObj, err := os.Open(orderNumbersFile)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("error during opening file: %v", err)
	}
	defer orderNumbersFileObj.Close()
	if !os.IsNotExist(err) {
		parser := json.NewDecoder(orderNumbersFileObj)
		if err = parser.Decode(&orderNumbers); err != nil {
			return nil, fmt.Errorf("error during decoding numbers db: %v", err)
		}
	}

	return orderNumbers, nil
}

func saveOrderNumbers(numbers map[string]interface{}) error {
	orderNumbersFileObj, err := os.OpenFile(orderNumbersFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("unable to open file for writing: %v", err)
	}
	defer orderNumbersFileObj.Close()
	encoder := json.NewEncoder(orderNumbersFileObj)
	return encoder.Encode(numbers)
}

func getNextInvoiceNumber(issueDate time.Time, contractorId string) (string, error) {
	var err error
	var exists bool
	var orderNumber interface{}
	var numbersPerMonth map[string]float64

	orderNumbers, err := parseOrderNumbers()
	if err != nil {
		return "", fmt.Errorf("cannot parse order numbers db: %v", err)
	}

	dateString := issueDate.Format("2006-01")

	numbersPerMonth, exists = orderNumbers[dateString].(map[string]float64)

	if !exists {
		// that's just a entry per month which may not exist
		numbersPerMonth = make(map[string]float64)
		numbersPerMonth[globalKey] = 1
		numbersPerMonth[contractorId] = 1
		orderNumbers[dateString] = numbersPerMonth
	}

	// fmt.Printf("orderNubmers: %+v\n", orderNumbers)
	// for key, value := range orderNumbers[dateString].(map[string]float64) {
	// 	numbersPerMonth[key] = value
	// }
	// numbersPerMonth = orderNumbers[dateString].(map[string]int)
	// let's check a number per contractor within month
	if _, exists = numbersPerMonth[contractorId]; !exists {
		// it does not, therefore let's bump the one from default entry
		if _, exists = numbersPerMonth[globalKey]; !exists {
			// there is no global key - let's recreate it.
			numbersPerMonth[globalKey] = 0
		}
		numbersPerMonth[globalKey] += 1
		numbersPerMonth[contractorId] = numbersPerMonth[globalKey]
		orderNumbers[dateString] = numbersPerMonth
	}

	orderNumber = numbersPerMonth[contractorId]

	orderNumberTemplate := pongo2.Must(pongo2.FromString(config.Config.OrderNumberFormat))

	if err = saveOrderNumbers(orderNumbers); err != nil {
		return "", fmt.Errorf("unable to save order numbers")
	}

	return orderNumberTemplate.Execute(pongo2.Context{
		"year":  issueDate.Year(),
		"month": issueDate.Month(),
		"no":    int(orderNumber.(float64)),
	})
}
