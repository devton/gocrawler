// Package crawler provides the core logic to apply
// the policies rules to scrape over a big data HTML folder
package crawler

import (
	"fmt"
	"os"
	"time"

	"github.com/devton/gocrawler/policies"
	"github.com/devton/gocrawler/xutils"
	"github.com/fatih/color"
)

var crawledDataFolder string

type ScrapedData struct {
	FilePath string                 `json:"filepath"`
	Fields   map[string]interface{} `json:"fields"`
}

// ScrapOver scrape using policies rules over given dataFolder
// and return a map with all scraped data and the total of scraped objects.
func ScrapOver(policiesFound []*policies.Policy, dataFolder string) (scrapedObjects map[string][]*ScrapedData, totalObjects int) {
	runningOn := make(chan map[string][]*ScrapedData, len(policiesFound))
	crawledDataFolder = dataFolder

	for _, policy := range policiesFound {
		go func(p *policies.Policy) {
			runningOn <- map[string][]*ScrapedData{
				p.FileName: ScrapeData(p),
			}
		}(policy)
	}

	return waitScrap(runningOn, policiesFound)
}

// ScrapeData parses the files found on rules of
// then given policy and return an array with ScrapedData struct
func ScrapeData(p *policies.Policy) []*ScrapedData {
	listFiles := p.GetAvaiableFilesFrom(crawledDataFolder)
	scrapedChannel := make(chan *ScrapedData, len(listFiles))

	for _, it := range listFiles {
		go func(item string, po *policies.Policy, c chan *ScrapedData) {
			b, _ := os.Open(item)
			defer b.Close()

			c <- &ScrapedData{
				FilePath: item,
				Fields:   p.ApplyRules(b),
			}
		}(it, p, scrapedChannel)
	}

	return waitAsyncScrape(scrapedChannel, len(listFiles))
}

// waitScrap waits all scraping data from given policies
// and returns a map of ScrapedData and total of objects scraped.
func waitScrap(channel chan map[string][]*ScrapedData, policiesFound []*policies.Policy) (scrapedObjects map[string][]*ScrapedData, totalObjects int) {
	responses := []string{}
	totalScrapedObjects := map[string][]*ScrapedData{}
	totalObjects = 0

	for {
		select {
		case result := <-channel:
			for key, data := range result {
				responses = append(responses, key)
				totalObjects += len(data)
				totalScrapedObjects[key] = append(totalScrapedObjects[key], data...)
			}

			if len(responses) == len(policiesFound) {
				return totalScrapedObjects, totalObjects
			}
		case <-time.After(30 * time.Millisecond):
		}
	}
}

// waitAsyncScrape waits for AsyncScrape ends
// and return a array with all ScrapedData
func waitAsyncScrape(scrapedChannel chan *ScrapedData, totalFiles int) []*ScrapedData {
	var totalScrapedObjects []*ScrapedData
	for {
		select {
		case asyncResult := <-scrapedChannel:
			fmt.Printf("%s", xutils.ColorSprint(color.FgGreen, "."))
			totalScrapedObjects = append(totalScrapedObjects, asyncResult)
			if len(totalScrapedObjects) == totalFiles {
				return totalScrapedObjects
			}
		}
	}
}
