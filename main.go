package main

import (
	"config"
	"discord"
	"os/signal"
	"syscall"
	"os"
	"restapi"
	"flag"
	"net/http"
)

func httpListen(keyPtr *string, certPtr *string, portPtr *string) {
	if *keyPtr != "" && *certPtr != "" {
		http.ListenAndServeTLS(":" + *portPtr, *certPtr, *keyPtr, nil)
	} else {
		http.ListenAndServe(":" + *portPtr, nil)
	}
}

func main() {
	keyPtr := flag.String("key", "", "path to key.pem file")
	certPtr := flag.String("cert", "", "path to cert.pem file")
	portPtr := flag.String("port", "9002", "port for http connections")

	flag.Parse()

	// set up rest api
	restapi.New()

	go httpListen(keyPtr, certPtr, portPtr)

	// setup main config
	config.Globalcfg = config.ReadConfig("./config.ini")
	dc := discord.NewDiscordBot(config.Globalcfg)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dc.Session.Close()
}
