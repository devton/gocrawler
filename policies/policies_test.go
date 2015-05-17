package policies

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindPolicies(t *testing.T) {
	assert := assert.New(t)
	policiesPath := "../test_files/policies"
	policiesFound := FindPolicies(policiesPath)
	policy := policiesFound[0]

	assert.Equal(1, len(policiesFound), "Should found a least one policy")
	assert.Equal("foo", policy.DomainFolder, "Should be filled policy.DomainFolder")
	assert.Equal([]string{"bar"}, policy.ResourcePaths, "Should be filled policy.ResourcePaths")
}
