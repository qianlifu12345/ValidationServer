package main

import (
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/rancher-auth-filter-service/manager"
	"github.com/rancher/rancher-auth-filter-service/service"
)

var VERSION = "v0.1.0-dev"

func main() {
	go manager.GetCommand(os.Args)
	log.Infof("Starting authantication filtering Service")
	router := service.NewRouter()
	http.Handle("/", router)
	// serverString := ":" + manager.Port
	serverString := ":8080"
	log.Fatal(http.ListenAndServe(serverString, nil))
}
