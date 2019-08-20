package common

import (
	"fmt"
	"os"
	"testing"
)

func Test_query(t *testing.T) {

	db, err := NewInstance("/Users/clark/workspace/golang/src/github.com/elastos/Elastos.ELA.Elephant.Node/elastos/data/ext")
	if err != nil {
		println(err.Error())
		os.Exit(-1)
	}
	height := "241762"
	l, err := db.Qu(`select m.* from (select ifnull(a.producer_public_key,b.ownerpublickey) as producer_public_key , ifnull(a.value,0) as value , b.* from
chain_producer_info b left join 
(select A.producer_public_key , ROUND(sum(value),8) as value from chain_vote_info A where (A.cancel_height > ` + height + ` or
cancel_height is null) and height <= ` + height + ` group by producer_public_key) a on a.producer_public_key = b.ownerpublickey
order by value desc) m `)
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Printf("%v\n", e.Value)
	}

}
