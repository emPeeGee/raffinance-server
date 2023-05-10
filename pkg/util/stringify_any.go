package util

import (
	"encoding/json"
)

func StringifyAny(object any) string {
	s, err := json.MarshalIndent(object, "", "\t")
	if err != nil {
		return "\nThe object cannot be prettified!\n"
	}

	return string(s)
}
