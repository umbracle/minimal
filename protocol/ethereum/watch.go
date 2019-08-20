package ethereum

import (
	"context"
	"fmt"

	"github.com/umbracle/minimal/types"
)

type pendingBlock struct {
	hash types.Hash
	num  uint64
	peer string
}

const (
	propagationRange = 10
)

type Watch struct {
	b       *Backend
	blockCh chan *pendingBlock
	list    []*pendingBlock
}

func (w *Watch) run() {
	// timer
	w.blockCh = make(chan *pendingBlock, 1)
	w.list = []*pendingBlock{}

	for {
		p := <-w.blockCh

		fmt.Println("_ DOWNLOAD _")
		fmt.Println(p.num)

		proto := w.b.getTarget(p.peer)
		header, err := proto.RequestHeaderByHashSync(context.Background(), p.hash)
		if err != nil {
			panic(err)
		}

		// localHeight, _ := w.b.blockchain.Header()

		if err := w.b.blockchain.WriteBlocks([]*types.Block{&types.Block{Header: header}}); err != nil {
			fmt.Println("-- err --")
			fmt.Println(err)
		}
	}
}

func (w *Watch) newHeaderUpdate(n uint64) {
	fmt.Println("____ HEADER UPDATE ___")
	// NOTE, this is called with the backend commitData but it does not ensure its the 'last'
	// data commited, will work for tests though

	size := len(w.list)
	if size == 0 {
		fmt.Println("-- back --")
		return
	}

	fmt.Println("-- vals --")
	fmt.Println(w.list[0].num)
	fmt.Println(n)

	if w.list[0].num >= n-1 {
		// do it!!
		fmt.Println("LLETS DO IT")
	}
}

func (w *Watch) isSynced() bool {
	// local height
	height, ok := w.b.blockchain.Header()
	if !ok {
		panic("header not found?")
	}

	pp := w.b.bestOne()

	return inRange(height.Number, pp.HeaderNumber, 10)
}

func (w *Watch) addPending(p *pendingBlock) {
	fmt.Println("-- ADD PENDING --")
	fmt.Println(p.num)

	height, _ := w.b.blockchain.Header()

	fmt.Println("-- height --")
	fmt.Println(height.Number)

	if height.Number >= p.num-1 {
		// sned it directly
		select {
		case w.blockCh <- p:
		default:
		}
		return
	}

	// relocate for later
	size := len(w.list)
	if size == 0 {
		w.list = append(w.list, p)
	} else {
		if w.list[size-1].num-1 == p.num {
			w.list = append(w.list, p)
		}
	}
}

// TODO, what happens if we get a notification with an invalid hash, that actually can
// happen at any time, three options: 1. its a microfork so we should have the parent
// 2. its an invalid message (it will fail), 3. notification from another long fork which we dont follow.

func (w *Watch) notify(peerID string, request *newBlockData) {

	num := request.Block.Number()
	b := request.TD.String()
	c := request.Block.Header.Difficulty

	fmt.Printf("===> NOTIFY (%s) Block: %d Difficulty %d. Total: %s\n", peerID, num, c, b)

	if !w.isSynced() {
		return
	}

	// local height
	height, ok := w.b.blockchain.Header()
	if !ok {
		panic("header not found?")
	}

	// height - 10 wont work if the number is less than 10 it will overflow
	if inRange(height.Number-10, height.Number, int64(num)) {
		// valid to download
		w.addPending(&pendingBlock{request.Block.Hash(), num, peerID})
	}
}

func (w *Watch) announce(peerID string, ann []*announcement) {
	// TODO, not used for now but it will be important later

	fmt.Printf("===> NOTIFY (%s) HASHES\n", peerID)
	fmt.Println(peerID, ann)

	if !w.isSynced() {
		return
	}

	// local height
	height, ok := w.b.blockchain.Header()
	if !ok {
		panic("header not found?")
	}

	// TODO, check the hash here too
	// Actually, the list will be a set of hash and number and we ask by hash which is easier

	for _, i := range ann {
		if inRange(height.Number-10, height.Number, int64(i.Number)) {
			// valid to download
		}
	}
}

func inRange(a, b uint64, r int64) bool {
	x := int64(a - b)
	if x < 0 {
		x = -1 * x
	}
	return x < r
}
