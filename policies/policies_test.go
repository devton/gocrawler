package policies

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindPolicies(t *testing.T) {
	assert := assert.New(t)
	policiesPath := "../test_files/policies"
	policiesFound := FindPolicies(policiesPath)
	policy := policiesFound[0]

	assert.Equal(1, len(policiesFound), "Should found a least one policy")
	assert.Equal("example.com", policy.DomainFolder, "Should be filled policy.DomainFolder")
	assert.Equal([]string{"path"}, policy.ResourcePaths, "Should be filled policy.ResourcePaths")
	assert.Equal(3, policy.MaxDepthForFindResource, "Should be filled policy.ResourcePaths")
}

func TestGetAvaiableFilesFrom(t *testing.T) {
	assert := assert.New(t)
	policy := Policy{
		DomainFolder:            "example.com",
		ResourcePaths:           []string{"path"},
		MaxDepthForFindResource: 3,
	}

	totalFiles := policy.GetAvaiableFilesFrom("../test_files/")
	assert.Equal(3, len(totalFiles), "should find all avaiable files")
}

func TestApplyRules(t *testing.T) {
	assert := assert.New(t)
	policiesPath := "../test_files/policies"
	policiesFound := FindPolicies(policiesPath)
	policy := policiesFound[0]

	files := policy.GetAvaiableFilesFrom("../test_files/")
	b, _ := os.Open(files[0])
	defer b.Close()
	parsedBody := policy.ApplyRules(b)

	assert.Equal("foo", reflect.ValueOf(parsedBody["text"]).String(), "should have filled field text")
	assert.Equal("bar", reflect.ValueOf(parsedBody["link"]).String(), "should have filled field link")
	assert.Equal(reflect.Slice, reflect.ValueOf(parsedBody["list"]).Kind(), "should have filled field list with string array")
	assert.Equal(2, reflect.ValueOf(parsedBody["list"]).Len(), "should have correct size the list field")
}

func TestFieldValue(t *testing.T) {
	assert := assert.New(t)
	policiesPath := "../test_files/policies"
	policiesFound := FindPolicies(policiesPath)
	policy := policiesFound[0]

	files := policy.GetAvaiableFilesFrom("../test_files/")
	b, _ := os.Open(files[0])
	defer b.Close()
	parsedBody := policy.ApplyRules(b)
	list := policy.FieldValue(parsedBody, "list").([]interface{})

	assert.Equal("foo", policy.FieldValue(parsedBody, "text"), "should have filled field text")
	assert.Equal("bar", policy.FieldValue(parsedBody, "link"), "should have filled field link")
	assert.Equal(2, len(list), "should have correct size the list field")
	assert.Equal("bar1", list[0], "should have correct size the list field")
	assert.Equal("bar2", list[1], "should have correct size the list field")
}
