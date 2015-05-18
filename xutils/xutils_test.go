package xutils

import (
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

func TestStringInSlice(t *testing.T) {
	assert := assert.New(t)
	slice := []string{"A", "B", "C"}

	assert.Equal(true, StringInSlice("C", slice), "Should be true when found string on slice")
	assert.Equal(false, StringInSlice("D", slice), "Should be false when can't found string on slice")
}

func TestExistsPath(t *testing.T) {
	assert := assert.New(t)
	exists, _ := ExistsPath("../test_files")
	invalid, _ := ExistsPath("../lorem_ipsum")

	assert.Equal(true, exists, "Should be true when path exists")
	assert.Equal(false, invalid, "Should be false when path does not exist")
}

func TestColorSprint(t *testing.T) {
	assert := assert.New(t)
	result := ColorSprint(color.FgGreen, "foo bar")

	assert.Equal("\x1b[32mfoo bar\x1b[0m", result, "Should return the text string with color")
}

func TestGetFilesOnDir(t *testing.T) {
	assert := assert.New(t)
	zeroMaxDepth := GetFilesOnDir("../test_files/example.com", 0, 0)
	oneMaxDepth := GetFilesOnDir("../test_files/example.com", 0, 1)

	assert.Equal(1, len(zeroMaxDepth), "Should return one file when run with 0 max depth")
	assert.Equal(2, len(oneMaxDepth), "Should return one file when run with 1 max depth")
}
