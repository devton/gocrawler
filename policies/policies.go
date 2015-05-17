package policies

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/devton/xporter/xutils"
	"github.com/fatih/color"
)

var policyLogTag = "[policy-finder]"
var policiesFound []string
var policiesSlice []*Policy

type Policy struct {
	FileName                string                       `json:"filename"`
	DomainFolder            string                       `json:"domain_folder"`
	ResourcePaths           []string                     `json:"resource_paths"`
	MaxDepthForFindResource int                          `json:"max_depth_to_find_resource"`
	Fields                  map[string]map[string]string `json:"fields"`
}

//Search by policy json files inside policyPath every 10s
func FindPolicies(policyPath string) []*Policy {
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
		return policiesSlice
	}
}

func PolicyWalk(path string, info os.FileInfo, err error) error {
	if info != nil && !info.IsDir() && filepath.Ext(path) == ".json" {
		fmt.Printf("%s", xutils.ColorSprint(color.FgMagenta, "."))

		body, err := ioutil.ReadFile(path)
		if err != nil {
			color.Red("can't read file %s -> %#v", filepath.Base(path), err)
		}

		var policy Policy

		dec := json.NewDecoder(bytes.NewReader(body))
		if err := dec.Decode(&policy); err != nil {
			color.Red("can't parse file %s -> %#v", filepath.Base(path), err)
		} else {
			if !xutils.StringInSlice(path, policiesFound) {
				policy.FileName = filepath.Base(path)
				policiesFound = append(policiesFound, path)
				policiesSlice = append(policiesSlice, &policy)
			}
		}
	}
	return nil
}
