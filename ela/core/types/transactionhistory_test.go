package types

import (
	"sort"
	"testing"
)

func TestSort(t *testing.T) {
	t1 := TransactionHistory{
		Height: 10,
	}
	t2 := TransactionHistory{
		Height: 2,
	}
	t0 := TransactionHistorySorter{
		t1,
		t2,
	}
	sort.Sort(t0)
}
