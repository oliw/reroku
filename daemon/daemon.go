package daemon

import (
	log "github.com/Sirupsen/logrus"
	"github.com/oliw/reroku/server"
)

func Start() error {
	log.Info("Starting daemon")
	srv, err := server.NewServer()
	if err != nil {
		return err
	}

	address := ""
	address += server.DEFAULTHTTPHOST
	address += ":"
	address += server.DEFAULTHTTPPORT
	err = server.ListenAndServe("tcp", address, srv)
	if err != nil {
		return err
	}
	return nil
}
