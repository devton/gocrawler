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
		fmt.Printf("%s check policy file %s\n",
			xutils.ColorSprint(color.FgMagenta, policyLogTag),
			xutils.ColorSprint(color.FgGreen, info.Name()))

		if !xutils.StringInSlice(path, policiesFound) {
			policiesFound = append(policiesFound, path)
		}
	}
	return nil
}
