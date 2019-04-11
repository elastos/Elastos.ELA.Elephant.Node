package common

import "github.com/elastos/Elastos.ELA/common"

func ContainsU168(c common.Uint168, s []common.Uint168) bool {
	for _, v := range s {
		if v == c {
			return true
		}
	}
	return false
}
