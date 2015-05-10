package crawler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/devton/xporter/policies"
	"github.com/devton/xporter/xutils"
	"github.com/fatih/color"
)

var crawlerLogTag = "[crawler]"
var crawledDataFolder string

func StartOver(policiesFound []string, dataFolder string) error {
	runningOn := make(chan string, len(policiesFound))
	responses := []string{}
	crawledDataFolder = dataFolder

	for _, policy := range policiesFound {
		fmt.Printf("%s runing for policy -> %s\n",
			xutils.ColorSprint(color.FgMagenta, crawlerLogTag),
			xutils.ColorSprint(color.FgGreen, filepath.Base(policy)))

		go AsyncCrawler(policy, runningOn)
	}

	for {
		select {
		case result := <-runningOn:
			responses = append(responses, result)
			if len(responses) == len(policiesFound) {
				fmt.Printf("\ndone all\n")
				return nil
			}
		case <-time.After(50 * time.Millisecond):
		}
	}
}

func AsyncCrawler(policyPath string, c chan string) {
	fmt.Printf("%s starting...\n", PolicyCrawlerLabel(policyPath))
	body, err := ioutil.ReadFile(policyPath)
	if err != nil {
		fmt.Errorf("%s can't read file %s...\n",
			PolicyCrawlerLabel(policyPath),
			xutils.ColorSprint(color.FgGreen, policyPath))
		c <- policyPath
	}

	var currentPolicy policies.Policy

	dec := json.NewDecoder(bytes.NewReader(body))
	if err := dec.Decode(&currentPolicy); err != nil {
		fmt.Errorf("%s can't parse policy file %s...\n",
			PolicyCrawlerLabel(policyPath),
			xutils.ColorSprint(color.FgGreen, policyPath))
		c <- policyPath
	}

	if err := SearchDataFor(&currentPolicy); err != nil {
		fmt.Errorf("%s search on crawled data folder for policy %s...\n",
			PolicyCrawlerLabel(policyPath),
			xutils.ColorSprint(color.FgGreen, policyPath))
		c <- policyPath
	}
}

func PolicyCrawlerLabel(policyPath string) string {
	return xutils.ColorSprint(color.FgCyan,
		fmt.Sprintf("[policy-crawler %s]", filepath.Base(policyPath)))
}

func SearchDataFor(p *policies.Policy) error {
	coursePath := path.Join(crawledDataFolder, p.DomainFolder)
	for _, cpath := range p.CoursePaths {
		files := GetFilesOnDir(path.Join(coursePath, cpath), 0, p.MaxDepthForCourse)

		for _, item := range files {
			fmt.Printf("%s scraping page file -> %s\n",
				xutils.ColorSprint(color.FgYellow, "[crawler-scraper]"),
				xutils.ColorSprint(color.FgGreen, item))

			b, _ := os.Open(item)
			queryDoc, _ := goquery.NewDocumentFromReader(b)

			for k, v := range p.Fields {
				fmt.Printf("%s -> %s\n", k, queryDoc.Find(v).Text())
			}
			b.Close()
		}
	}

	return nil
}

func GetFilesOnDir(filePath string, depth int, maxDepth int) []string {
	tempFiles := []string{}

	if depth <= maxDepth {
		files, _ := ioutil.ReadDir(filePath)
		for _, item := range files {
			if item.IsDir() {
				for _, recItem := range GetFilesOnDir(
					path.Join(filePath, item.Name()), (depth + 1), maxDepth) {
					tempFiles = append(tempFiles, recItem)
				}
			} else {
				tempFiles = append(tempFiles, path.Join(filePath, item.Name()))
			}
		}
	}

	return tempFiles
}
