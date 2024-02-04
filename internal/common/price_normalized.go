package common

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type PriceNormalized struct {
	Price      int
	Multiplier int
}

func ParsePrice(input string) (PriceNormalized, error) {
	var output PriceNormalized
	var err error

	dotPosition := strings.Index(input, ".")
	if dotPosition == -1 {
		input = input + ".0"
		dotPosition = len(input) - 2
	}
	output.Multiplier = len(input) - dotPosition - 1
	// now that we have the dot location we can get rid of it and parse the remaining
	// string as int
	if output.Price, err = strconv.Atoi(strings.Replace(input, ".", "", 1)); err != nil {
		return output, err
	}
	return output, nil
}

// https://stackoverflow.com/questions/13020308/how-to-fmt-printf-an-integer-with-thousands-comma
func (p PriceNormalized) Format(thousandSeparator string) string {
	quantizer := int(math.Pow10(p.Multiplier))
	remainder := p.Price % quantizer
	output := strconv.Itoa(p.Price / quantizer)
	startOffset := 3
	if p.Price < 0 {
		startOffset++
	}
	for outputIndex := len(output); outputIndex > startOffset; {
		outputIndex -= 3
		output = output[:outputIndex] + thousandSeparator + output[outputIndex:]
	}
	if remainder > 0 {
		return fmt.Sprintf("%s.%02d", output, remainder)
	}
	return output
}

func (p PriceNormalized) AsFloat() float64 {
	return float64(p.Price) / math.Pow10(p.Multiplier)
}
