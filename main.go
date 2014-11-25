package main

import (
	"github.com/benschw/go-protect/protect"
	// "gopkg.in/yaml.v1"
	"flag"
	"fmt"
	"github.com/goraft/raft"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var verbose bool
	var trace bool
	var debug bool
	var bootstrap bool
	var dataDir string
	var raftAddr string
	var join string
	var apiAddr string

	flag.BoolVar(&verbose, "v", false, "verbose logging")
	flag.BoolVar(&trace, "trace", false, "Raft trace debugging")
	flag.BoolVar(&debug, "debug", false, "Raft debugging")

	flag.BoolVar(&bootstrap, "bootstrap", false, "Bootstrap Raft cluster")
	flag.StringVar(&dataDir, "data", "./data", "Raft data dir path")
	flag.StringVar(&raftAddr, "raft", "localhost:5000", "Raft hostname:port")
	flag.StringVar(&join, "join", "", "Leader host:port to join")
	flag.StringVar(&apiAddr, "api", "localhost:6000", "API hostname:port")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [arguments] <command> \n", os.Args[0])
		flag.PrintDefaults()
	}

	log.SetFlags(0)
	flag.Parse()
	if verbose {
		log.Print("Verbose logging enabled.")
	}
	if trace {
		raft.SetLogLevel(raft.Trace)
		log.Print("Raft trace debugging enabled.")
	} else if debug {
		raft.SetLogLevel(raft.Debug)
		log.Print("Raft debugging enabled.")
	}

	log.SetFlags(log.LstdFlags)

	// pull desired command/operation from args
	if flag.NArg() == 0 {
		flag.Usage()
		log.Fatal("Command argument required")
	}
	command := flag.Arg(0)

	// Populate Config
	rAddrParts := strings.Split(raftAddr, ":")
	rPort, _ := strconv.Atoi(rAddrParts[1])
	aAddrParts := strings.Split(apiAddr, ":")
	aPort, _ := strconv.Atoi(aAddrParts[1])
	cfg := protect.Config{
		RaftHost:  rAddrParts[0],
		RaftPort:  rPort,
		ApiHost:   aAddrParts[0],
		ApiPort:   aPort,
		DataDir:   dataDir,
		JoinAddr:  join,
		Bootstrap: bootstrap,
	}

	// Run Main App
	switch command {
	case "serve":
		svc := protect.Service{}

		if err := svc.Run(cfg); err != nil {
			log.Fatal(err)
		}
	default:
		flag.Usage()
		log.Fatalf("Unknown Command: %s", command)
	}

	// app := cli.NewApp()
	// app.Name = "go-protect"
	// app.Usage = "work with the `go-protect` service"
	// app.Version = "0.0.1"

	// app.Flags = []cli.Flag{
	// 	cli.StringFlag{"config, c", "config.yaml", "config file to use"},
	// }

	// app.Commands = []cli.Command{
	// 	{
	// 		Name:  "server",
	// 		Usage: "Run the http server",
	// 		Action: func(c *cli.Context) {
	// 			cfg, err := getConfig(c)
	// 			if err != nil {
	// 				log.Fatal(err)
	// 				return
	// 			}

	// 			svc := protect.Service{}

	// 			if err = svc.Run(cfg); err != nil {
	// 				log.Fatal(err)
	// 			}
	// 		},
	// 	},
	// 	// {
	// 	// 	Name:  "migratedb",
	// 	// 	Usage: "Perform database migrations",
	// 	// 	Action: func(c *cli.Context) {
	// 	// 		cfg, err := getConfig(c)
	// 	// 		if err != nil {
	// 	// 			log.Fatal(err)
	// 	// 			return
	// 	// 		}

	// 	// 		svc := service.TodoService{}

	// 	// 		if err = svc.Migrate(cfg); err != nil {
	// 	// 			log.Fatal(err)
	// 	// 		}
	// 	// 	},
	// 	// },
	// }
	// app.Run(os.Args)

}
