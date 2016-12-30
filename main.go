package main

import (
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/rancher-auth-filter-service/service"
	"github.com/urfave/cli"
)

var VERSION = "v0.1.0-dev"
var url = "http://54.255.182.226"
var port = "8080"

func main() {
	log.Infof("Starting authantication filtering Service")
	app := cli.NewApp()
	app.Name = "rancher-auth-filter-service"
	app.Version = "v0.1.0-dev"
	app.Usage = "Rancher authantication Filter Service"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "rancherUrl",
			Value: "http://54.255.182.226",
			Usage: "Rancher server url",
		},
		cli.StringFlag{
			Name:  "localport",
			Value: "8080",
			Usage: "Local server port ",
		},
	}

	app.Action = func(c *cli.Context) error {
		url = c.String("url")
		port = c.String("port")
		return nil
	}

	app.Run(os.Args)

	router := service.NewRouter()
	http.Handle("/", router)
	serverString := ":" + port
	log.Fatal(http.ListenAndServe(serverString, nil))
}
