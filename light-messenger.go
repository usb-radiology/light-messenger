package main

import (
	"log"

	"github.com/usb-radiology/light-messenger/src/configuration"
	"github.com/usb-radiology/light-messenger/src/server"
	"github.com/usb-radiology/light-messenger/src/version"
)

// TODO: add command parameters to enable script-like workflow

func main() {
	log.Printf("%s %s", version.Version, version.BuildTime)

	initConfig, err := configuration.LoadAndSetConfiguration("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	httpServer := server.InitServer(initConfig)
	server.Start(httpServer, initConfig.Server.HTTPPort)

}
