package category

import (
	"fmt"
	"strings"
)

func checkForBlacklist(name string) error {
	for _, blacklistedName := range CategoryBlacklist {
		if strings.EqualFold(name, blacklistedName) {
			return fmt.Errorf("category name %s is not allowed", name)
		}
	}

	return nil
}
