package main

import (
	"config"
	"discord"
	"os/signal"
	"syscall"
	"os"
)

func main() {
	// setup main config
	config.Globalcfg = config.ReadConfig("./config.ini")
	dc := discord.NewDiscordBot(config.Globalcfg)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	
	dc.Session.Close()
}
