package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/codegangsta/cli"
	"github.com/devton/gocrawler/crawler"
	"github.com/devton/gocrawler/policies"
	"github.com/devton/gocrawler/xutils"
	"github.com/fatih/color"
)

const Version string = "0.0.0"

func main() {
	app := cli.NewApp()
	app.Authors = app.Authors[:0]
	app.Name = "gocrawler"
	app.Version = Version
	app.Usage = ""
	app.Author = "Antônio Roberto"
	app.Email = "forevertonny@gmail.com"
	app.Commands = []cli.Command{
		cli.Command{
			Name:  "craw",
			Usage: "Start crawling and exporting data from crawled site documents",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "crawled-files, cf",
					Value: "/path/to/.crawled_data/",
					Usage: "folder with downloaded HTML pages",
				},
				cli.StringFlag{
					Name:  "policies-path, pf",
					Value: "/path/to/.crawler-json-policies/",
					Usage: "folder with json policies to extract data",
				},
				cli.StringFlag{
					Name:  "output-files, of",
					Value: "gocrawler-output",
					Usage: "default output json files",
				},
			},
			Action: func(c *cli.Context) {
				color.Green("starting crawler...")
				for {
					//Find policies
					policiesFound := policies.FindPolicies(c.String("policies-path"))

					startCrawlingAt := time.Now()
					//Scrap data from files using rules descibres at policies found
					scrapedData, totalObjects := crawler.ScrapOver(policiesFound, c.String("crawled-files"))
					color.Magenta("\nTotal files parsed -> %d", totalObjects)

					//Save scraped data into .json files
					existsPath, _ := xutils.ExistsPath(c.String("output-files"))
					if !existsPath {
						os.MkdirAll(c.String("output-files"), 0755)
					}

					for k, v := range scrapedData {
						color.Yellow("saving json file for policy -> %s", path.Join(c.String("output-files"), filepath.Base(k)))
						jsonData, _ := json.Marshal(v)
						ioutil.WriteFile(path.Join(c.String("output-files"), filepath.Base(k)), jsonData, 0755)
					}

					elapsedCrawlingTime := time.Since(startCrawlingAt)
					color.Green("scraping total time took %s", elapsedCrawlingTime)
					select {
					case <-time.After(30 * time.Minute):
						color.Green("rerunning crawler... ;)")
					}
				}
			},
		},
	}

	app.Run(os.Args)
}
