package itrie

import (
	"fmt"
	"testing"

	"github.com/umbracle/minimal/helper/hex"
	"github.com/umbracle/minimal/types"
)

func TestTrieExpansion(t *testing.T) {
	storage := NewMemoryStorage()
	batch := storage.Batch()

	state := NewState(storage)
	trie := state.NewSnapshot().(*Trie)

	data := hex.MustDecodeHex("0x123251456987548525698451254587452222222222222222222222")

	txn := trie.Txn()
	txn.batch = batch

	txn.Insert([]byte{0x1, 0x2}, data)
	txn.Insert([]byte{0x3, 0x4}, data)
	txn.Insert([]byte{0x5, 0x6}, data)

	txn.Show()

	root, err := txn.Hash()
	if err != nil {
		panic(err)
	}

	fmt.Println(root)
	txn.Commit()

	batch.Write()

	// fmt.Println(storage)

	state2 := NewState(storage)
	aux, _ := state2.NewSnapshotAt(types.BytesToHash(root))
	trie2 := aux.(*Trie)

	txn2 := trie2.Txn()

	fmt.Println("-- result --")
	fmt.Println(txn2.Lookup([]byte{0x1, 0x2}))

	txn2.Show()

}
