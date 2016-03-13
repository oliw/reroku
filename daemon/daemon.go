package daemon

import (
	"fmt"
	"time"
)

func Start() {
	for x := range time.Tick(3 * time.Second) {
		fmt.Printf("Daemon is running %d\n", x)
	}
}
