package state

import "github.com/elastos/Elastos.ELA/dpos/state"

var DefaultArbitratorsEx ArbitratorsEx

type ArbitratorsEx interface {
	state.Arbitrators
	ForceChange(height uint32) error
}
