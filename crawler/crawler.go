package crawler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/devton/xporter/policies"
	"github.com/devton/xporter/xutils"
	"github.com/fatih/color"
)

var crawlerLogTag = "[crawler]"
var crawledDataFolder string

type ScrapedData struct {
	Fields map[string]string `json:"fields"`
}

func ScrapOver(policiesFound []string, dataFolder string) (scrapedObjects map[string][]*ScrapedData, totalObjects int) {
	runningOn := make(chan map[string][]*ScrapedData, len(policiesFound))
	responses := []string{}
	totalScrapedObjects := map[string][]*ScrapedData{}
	totalObjects = 0
	crawledDataFolder = dataFolder

	fmt.Print("crawling for policies -> ")
	for _, policy := range policiesFound {
		fmt.Printf(" %s", xutils.ColorSprint(color.FgMagenta, filepath.Base(policy)))
		go AsyncCrawler(policy, runningOn)
	}

	for {
		select {
		case result := <-runningOn:
			for key, data := range result {
				responses = append(responses, key)
				totalObjects += len(data)
				for _, item := range data {
					totalScrapedObjects[key] = append(totalScrapedObjects[key], item)
				}
			}

			if len(responses) == len(policiesFound) {
				return totalScrapedObjects, totalObjects
			}
		case <-time.After(50 * time.Millisecond):
		}
	}

	return totalScrapedObjects, totalObjects
}

func AsyncCrawler(policyPath string, c chan map[string][]*ScrapedData) {
	fmt.Printf("%s", xutils.ColorSprint(color.FgCyan, "."))
	body, err := ioutil.ReadFile(policyPath)
	if err != nil {
		fmt.Errorf("can't read file %s...\n", policyPath)
	}

	var currentPolicy policies.Policy

	dec := json.NewDecoder(bytes.NewReader(body))
	if err := dec.Decode(&currentPolicy); err != nil {
		fmt.Errorf("can't parse policy file %s...\n", policyPath)
	}

	data := ScrapeData(&currentPolicy)
	c <- map[string][]*ScrapedData{
		policyPath: data,
	}
}

func ScrapeData(p *policies.Policy) []*ScrapedData {
	coursePath := path.Join(crawledDataFolder, p.DomainFolder)
	totalScrapedObjects := []*ScrapedData{}
	totalFiles := 0

	for _, cpath := range p.CoursePaths {
		files := GetFilesOnDir(path.Join(coursePath, cpath), 0, p.MaxDepthForFindResource)
		scrapedChannel := make(chan *ScrapedData, len(files))
		totalFiles += len(files)

		for _, it := range files {
			go AsyncScrape(it, p, scrapedChannel)
		}

		for {
			select {
			case asyncResult := <-scrapedChannel:
				fmt.Printf("%s", xutils.ColorSprint(color.FgGreen, "."))
				totalScrapedObjects = append(totalScrapedObjects, asyncResult)
				if len(totalScrapedObjects) == len(files) {
					return totalScrapedObjects
				}
			}
		}
	}

	return totalScrapedObjects
}

func AsyncScrape(item string, p *policies.Policy, c chan *ScrapedData) error {
	b, _ := os.Open(item)
	queryDoc, _ := goquery.NewDocumentFromReader(b)
	var body = map[string]string{}

	for k, v := range p.Fields {
		selector := queryDoc.Find(v["selector"])

		if filterString, ok := v["filters"]; ok {
			filter := strings.Split(filterString, ".")
			if len(filter) > 1 {
				switch filter[0] {
				case "attr":
					body[k] = selector.AttrOr(filter[1], "")
				}
			}

		} else {
			body[k] = selector.Text()
		}
	}

	b.Close()
	c <- &ScrapedData{Fields: body}
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
