package xutils

import (
	"os"

	"github.com/fatih/color"
)

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func ColorSprint(c color.Attribute, msg string) string {
	green := color.New(c).SprintFunc()
	return green(msg)
}

func ExistsPath(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
