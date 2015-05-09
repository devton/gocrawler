package main

import (
	"os"

	"github.com/codegangsta/cli"
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
					Value: "~/crawled_data/",
					Usage: "folder with downloaded HTML pages",
				},
				cli.StringFlag{
					Name:  "policies-files, pf",
					Value: "~/crawler-json-policies/",
					Usage: "folder with json policies to extract data",
				},
				cli.StringFlag{
					Name:  "elastic-url, eurl",
					Value: "http://127.0.0.1:9200",
					Usage: "elastic search host url",
				},
				cli.StringFlag{
					Name:  "default-es-index",
					Value: "xported-documents",
					Usage: "default elastic search index to save extracted info",
				},
			},
			Action: func(c *cli.Context) {
				color.Green("starting crawler...")
			},
		},
	}

	app.Run(os.Args)
}
