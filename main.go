package main

import (
	"github.com/elastos/Elastos.ELA.Utility/signal"
)

func main() {
	var interrupt = signal.NewInterrupt()
	defer func() {
		<-interrupt.C
	}()
	go Go(Version, GoVersion)
}
