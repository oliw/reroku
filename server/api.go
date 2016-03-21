package server

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"net"
	"net/http"
)

const DEFAULTHTTPHOST string = "127.0.0.1"
const DEFAULTHTTPPORT string = "2514"

func getVersion(srv *Server, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	v := srv.Version()
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	writeJSON(w, b)
	return nil
}

func writeJSON(w http.ResponseWriter, b []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func createRouter(srv *Server) (*mux.Router, error) {
	r := mux.NewRouter()
	m := map[string]map[string]func(*Server, http.ResponseWriter, *http.Request, map[string]string) error{
		"GET": {
			"/version": getVersion,
		},
	}

	for method, routes := range m {
		for route, fct := range routes {
			log.Debug("Registering %s %s", method, route)
			f := func(w http.ResponseWriter, r *http.Request) {
				log.Debug("Calling %s %s", method, route)
				if err := fct(srv, w, r, mux.Vars(r)); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
			r.Path(route).Methods(method).HandlerFunc(f)
		}
	}
	return r, nil
}

func ListenAndServe(protocol, address string, srv *Server) error {
	log.Info(fmt.Sprintf("Starting API at %s", address))
	r, err := createRouter(srv)
	if err != nil {
		return err
	}
	l, err := net.Listen(protocol, address)
	if err != nil {
		return err
	}
	httpSrv := http.Server{Addr: address, Handler: r}
	return httpSrv.Serve(l)
}
