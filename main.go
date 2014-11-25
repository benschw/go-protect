package main

import (
	"errors"
	"github.com/benschw/go-protect/protect"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"os"
)

func getConfig(c *cli.Context) (protect.Config, error) {
	yamlPath := c.GlobalString("config")
	config := protect.Config{}

	if _, err := os.Stat(yamlPath); err != nil {
		return config, errors.New("config path not valid")
	}

	ymlData, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal([]byte(ymlData), &config)
	return config, err
}

func main() {

	app := cli.NewApp()
	app.Name = "go-protect"
	app.Usage = "work with the `go-protect` service"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{"config, c", "config.yaml", "config file to use"},
	}

	app.Commands = []cli.Command{
		{
			Name:  "server",
			Usage: "Run the http server",
			Action: func(c *cli.Context) {
				cfg, err := getConfig(c)
				if err != nil {
					log.Fatal(err)
					return
				}

				svc := protect.Service{}

				if err = svc.Run(cfg); err != nil {
					log.Fatal(err)
				}
			},
		},
		// {
		// 	Name:  "migratedb",
		// 	Usage: "Perform database migrations",
		// 	Action: func(c *cli.Context) {
		// 		cfg, err := getConfig(c)
		// 		if err != nil {
		// 			log.Fatal(err)
		// 			return
		// 		}

		// 		svc := service.TodoService{}

		// 		if err = svc.Migrate(cfg); err != nil {
		// 			log.Fatal(err)
		// 		}
		// 	},
		// },
	}
	app.Run(os.Args)

}
