package dht_test

import (
	"testing"

	"github.com/wjh/zif/libzif/dht"
)

func TestNetDBInsert(t *testing.T) {
	// would be so cool if you had that public key for real :O
	addr := dht.NewAddress([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
		15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
		27, 28, 29, 30, 31, 32})

	addr2 := dht.NewAddress([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
		15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
		27, 28, 29, 30, 31, 33})

	db := dht.NewNetDB(addr)

	insert := func(a dht.Address, expectedLen int) {
		err := db.Insert(dht.NewKeyValue(a, []byte{13, 37}))

		if err != nil {
			t.Error(err.Error())
		}

		if db.TableLen() != expectedLen {
			t.Errorf("TableLen not correct: %d", db.TableLen())
		}
	}

	// two inserts should result in just one
	insert(addr, 1)
	insert(addr, 1)

	insert(addr2, 2)
}
