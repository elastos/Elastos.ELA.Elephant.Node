package main

import (
	"github.com/elastos/Elastos.ELA.Elephant.Node/ela"
	"github.com/elastos/Elastos.ELA.Utility/signal"
)

func main() {
	var interrupt = signal.NewInterrupt()
	go ela.Go()
	//go id.Go()
	<-interrupt.C
}
