package main

import (
	"log"

	"github.com/usb-radiology/light-messenger/src/configuration"
	"github.com/usb-radiology/light-messenger/src/server"
)

func main() {

	initConfig, err := configuration.LoadAndSetConfiguration("./config.json")
	if err != nil {
		log.Fatal(err)
	}

	httpServer := server.InitServer(initConfig)
	server.Start(httpServer, initConfig.Server.HTTPPort)

}
