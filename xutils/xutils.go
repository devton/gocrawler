package xutils

import "github.com/fatih/color"

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
