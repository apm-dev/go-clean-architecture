package main

import (
	"fmt"
	"github.com/apm-dev/go-clean-architecture/app"
	"os"
	"os/signal"
)

func main() {
	fmt.Println("** Starting Application **")
	app.StartApplication()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	// Block until a signal is received
	<-ch

	fmt.Println("** Stopping Application **")
	app.StopApplication()
}
