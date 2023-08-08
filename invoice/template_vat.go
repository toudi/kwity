package invoice

import (
	"fmt"
	"strconv"
)

type Vat struct {
	Rate        float32 `json:"rate"`
	Description string  `json:"label"`
}

func (v *Vat) GetDescription() string {
	if v.Description != "" {
		return v.Description
	}
	return fmt.Sprintf("%s %%", strconv.FormatFloat(float64(v.Rate*100), 'f', -1, 64))
}
