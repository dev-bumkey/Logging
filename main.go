package main

import (
	"os"
	"os/signal"

	"github.com/cocktailcloud/acloud-alarm-collector/cmd"
)

func main() {
	go cmd.Execute()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	cmd.Close()
	os.Exit(0)
}
