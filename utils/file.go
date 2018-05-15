package utils

import (
	"fmt"
	"os"
)

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("[FileExists] error(%s)", err)
			return false
		}
	}
	return true
}
