package main

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/rancher/rancher-auth-filter-service/manager"
	"github.com/rancher/rancher-auth-filter-service/service"
)

var VERSION = "v0.1.0-dev"

func main() {
	manager.GetCommand()
	logrus.Infof("Starting authantication filtering Service")
	router := service.NewRouter()
	http.Handle("/", router)
	serverString := ":" + manager.Port
	logrus.Fatal(http.ListenAndServe(serverString, nil))
}
