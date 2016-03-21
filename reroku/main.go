package main

import (
	"flag"
	"fmt"
	"github.com/oliw/reroku/client"
	"github.com/oliw/reroku/daemon"
	"log"
	"os"
)

func main() {
	flDaemon := flag.Bool("d", false, "Daemon mode")
	flag.Parse()
	if *flDaemon {
		fmt.Printf("Launching daemon\n")
		daemon.Start()
	} else {
		if err := client.ParseCommands(flag.Args()...); err != nil {
			log.Fatal(err)
			os.Exit(-1)
		}
	}
	return
}
