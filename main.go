package main

import (
	"github.com/benschw/go-protect/protect"
	"github.com/benschw/go-protect/raft/command"
	// "gopkg.in/yaml.v1"
	"flag"
	"fmt"
	"github.com/goraft/raft"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
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
	cmd := flag.Arg(0)

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

	rand.Seed(time.Now().UnixNano())

	raft.RegisterCommand(&command.WriteCommand{})

	// Run Main App
	switch cmd {
	case "serve":
		svc := protect.Service{}

		if err := svc.Run(cfg); err != nil {
			log.Fatal(err)
		}
	default:
		flag.Usage()
		log.Fatalf("Unknown Command: %s", cmd)
	}
}
