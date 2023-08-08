package invoice

import (
	"encoding/json"
	"os"
)

func (i *Invoice) renderJSON() error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(i)
}
