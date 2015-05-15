package policies

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/devton/xporter/xutils"
	"github.com/fatih/color"
)

var policyLogTag = "[policy-finder]"
var policiesFound []string

type Policy struct {
	DomainFolder            string                       `json:"domain_folder"`
	CoursePaths             []string                     `json:"course_paths"`
	MaxDepthForFindResource int                          `json:"max_depth_to_find_resource"`
	Fields                  map[string]map[string]string `json:"fields"`
}

//Search by policy json files inside policyPath every 10s
func FindPolicies(policyPath string) []string {
	fmt.Printf("%s looking for policies at %s\n",
		xutils.ColorSprint(color.FgMagenta, policyLogTag),
		xutils.ColorSprint(color.FgGreen, policyPath))

	filepath.Walk(policyPath, PolicyWalk)

	select {
	case <-time.After(500 * time.Millisecond):
		if len(policiesFound) == 0 {
			fmt.Printf("%s %s closing program...",
				xutils.ColorSprint(color.FgMagenta, policyLogTag),
				xutils.ColorSprint(color.FgRed, "no policies found."))
			os.Exit(1)
		}
		return policiesFound
	}
}

func PolicyWalk(path string, info os.FileInfo, err error) error {
	if info != nil && !info.IsDir() && filepath.Ext(path) == ".json" {
		fmt.Printf("%s", xutils.ColorSprint(color.FgMagenta, "."))

		if !xutils.StringInSlice(path, policiesFound) {
			policiesFound = append(policiesFound, path)
		}
	}
	return nil
}
