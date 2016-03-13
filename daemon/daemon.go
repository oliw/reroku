package daemon

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"time"
)

func Start() {
	log.Info("Starting daemon")
	for x := range time.Tick(3 * time.Second) {
		fmt.Printf("Daemon is running %d\n", x)
	}
}
