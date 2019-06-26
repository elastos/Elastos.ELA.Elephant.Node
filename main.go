package main

import (
	"github.com/elastos/Elastos.ELA.Elephant.Node/ela"
	"github.com/elastos/Elastos.ELA.Utility/signal"
)

var (
	// Build version generated when build program.
	Version string

	// The go source code version at build.
	GoVersion string
)

func main() {
	var interrupt = signal.NewInterrupt()
	defer func() {
		<-interrupt.C
	}()
	go ela.Go(Version, GoVersion)
	//go id.Go(Version, GoVersion)
}
