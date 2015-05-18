package xutils

import (
	"io/ioutil"
	"os"
	"path"

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

func GetFilesOnDir(filePath string, depth int, maxDepth int) []string {
	tempFiles := []string{}

	if depth <= maxDepth {
		files, _ := ioutil.ReadDir(filePath)
		for _, item := range files {
			if item.IsDir() {
				subFiles := GetFilesOnDir(
					path.Join(filePath, item.Name()), (depth + 1), maxDepth)

				for _, recItem := range subFiles {
					tempFiles = append(tempFiles, recItem)
				}
			} else {
				tempFiles = append(tempFiles, path.Join(filePath, item.Name()))
			}
		}
	}

	return tempFiles
}
