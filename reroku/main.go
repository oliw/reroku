package main

import (
	"flag"
	"fmt"
	"github.com/oliw/reroku/daemon"
)

func main() {
	flDaemon := flag.Bool("d", false, "Daemon mode")
	flag.Parse()
	if *flDaemon {
		fmt.Printf("Launching daemon\n")
		daemon.Start()
	}
	return
}