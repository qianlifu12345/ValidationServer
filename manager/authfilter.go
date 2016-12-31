package manager

import (
	"os"

	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var Url = "http://54.255.182.226:8080/"
var Port = "8080"

//GetCommand is to get the rancher server url and local port
func GetCommand() {
	app := cli.NewApp()
	app.Name = "rancher-auth-filter-service"
	app.Version = "v0.1.0-dev"
	app.Usage = "Rancher authantication Filter Service"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "rancherUrl",
			Value: "http://54.255.182.226:8080/",
			Usage: "Rancher server url",
		},
		cli.StringFlag{
			Name:  "localport",
			Value: "8080",
			Usage: "Local server port ",
		},
	}
	app.Action = func(c *cli.Context) error {
		Url = c.String("rancherUrl")
		Port = c.String("localport")
		logrus.Infof("URL:" + Url + "LocalPort:" + Port)
		// logrus.Infof(c.NArg)
		fmt.Println(c.NArg())
		return nil
	}

	app.Run(os.Args)

}
