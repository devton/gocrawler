package crawler

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/devton/xporter/policies"
	"github.com/devton/xporter/xutils"
	"github.com/fatih/color"
)

var crawledDataFolder string

type ScrapedData struct {
	FilePath string            `json:"filepath"`
	Fields   map[string]string `json:"fields"`
}

func ScrapOver(policiesFound []*policies.Policy, dataFolder string) (scrapedObjects map[string][]*ScrapedData, totalObjects int) {
	runningOn := make(chan map[string][]*ScrapedData, len(policiesFound))
	responses := []string{}
	totalScrapedObjects := map[string][]*ScrapedData{}
	totalObjects = 0
	crawledDataFolder = dataFolder

	fmt.Print("crawling for policies -> ")
	for _, policy := range policiesFound {
		fmt.Printf(" %s", xutils.ColorSprint(color.FgMagenta, policy.FileName))
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

func AsyncCrawler(currentPolicy *policies.Policy, c chan map[string][]*ScrapedData) {
	c <- map[string][]*ScrapedData{
		currentPolicy.FileName: ScrapeData(currentPolicy),
	}
}

func ScrapeData(p *policies.Policy) []*ScrapedData {
	coursePath := path.Join(crawledDataFolder, p.DomainFolder)
	totalScrapedObjects := []*ScrapedData{}
	totalFiles := 0

	for _, cpath := range p.ResourcePaths {
		files := xutils.GetFilesOnDir(path.Join(coursePath, cpath), 0, p.MaxDepthForFindResource)
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
	c <- &ScrapedData{
		FilePath: item,
		Fields:   body,
	}
	return nil
}
