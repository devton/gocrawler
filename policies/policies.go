// Package policies provides core logic to
// find and parse json policy to struct Policy
// and extract content using the defined policy rules
package policies

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/devton/gocrawler/xutils"
	"github.com/fatih/color"
)

var policiesFound []string
var policiesSlice []*Policy

// Policy struct responsibl to handle with
// json data from policy file
type Policy struct {
	FileName                string                       `json:"filename"`
	DomainFolder            string                       `json:"domain_folder"`
	ResourcePaths           []string                     `json:"resource_paths"`
	MaxDepthForFindResource int                          `json:"max_depth_to_find_resource"`
	Fields                  map[string]map[string]string `json:"fields"`
}

// GetAvaiableFilesFrom get all avaiables files from
// the given folder over the current Policy
func (p *Policy) GetAvaiableFilesFrom(folder string) []string {
	var files []string

	for _, cpath := range p.ResourcePaths {
		files = append(files, xutils.GetFilesOnDir(
			path.Join(folder, p.DomainFolder, cpath), 0, p.MaxDepthForFindResource)...)
	}

	return files
}

// ApplyRules apply policy fields rules over
// the given File
func (p *Policy) ApplyRules(b *os.File) map[string]interface{} {
	body := make(map[string]interface{})
	queryDoc, _ := goquery.NewDocumentFromReader(b)

	for k, v := range p.Fields {
		selector := queryDoc.Find(v["selector"])

		if filterString, ok := v["filters"]; ok {
			filter := strings.Split(filterString, "|")
			if len(filter) > 1 {
				switch filter[0] {
				case "map":
					body[k] = selector.Map(func(i int, s *goquery.Selection) string {
						return smallFilterApply(s, filter[1])
					})
				}
			} else {
				body[k] = smallFilterApply(selector, filterString)
			}
		} else {
			body[k] = selector.Text()
		}
	}

	return body
}

// FieldValue returns the correct value inside an interface
// for given fieldName over Parsed map string
func (p *Policy) FieldValue(virtualParsed map[string]interface{}, fieldName string) interface{} {
	val := reflect.ValueOf(virtualParsed[fieldName])

	switch val.Kind() {
	case reflect.Slice:
		ret := make([]interface{}, val.Len())
		for i := 0; i < val.Len(); i++ {
			ret[i] = val.Index(i).Interface()
		}
		return ret
	default:
		return val.Interface()
	}
}

// smallFilterApply apply small filter on
// goquery selection using the given filterString
func smallFilterApply(s *goquery.Selection, filterString string) string {
	filter := strings.Split(filterString, ".")

	if len(filter) > 1 {
		switch filter[0] {
		case "attr":
			return s.AttrOr(filter[1], "")
		default:
			return ""
		}
	}

	return ""
}

// FindPolicies search by valid policies inside a given path
// and return an array of Policy struct.
func FindPolicies(policyPath string) []*Policy {
	fmt.Printf("looking for policies at %s\n",
		xutils.ColorSprint(color.FgGreen, policyPath))

	filepath.Walk(policyPath, PolicyWalk)
	return waitWalk()
}

// waitWalk waits for filepath.Walk return
// a least one policy in 500ms
func waitWalk() []*Policy {
	select {
	case <-time.After(500 * time.Millisecond):
		if len(policiesFound) == 0 {
			fmt.Printf("%s bleh...",
				xutils.ColorSprint(color.FgRed, "no policies found."))
		}
		return policiesSlice
	}
}

// PolicyWalk is a walk func for filepath.Walk, parse a json policy
// file into a Policy struct
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
