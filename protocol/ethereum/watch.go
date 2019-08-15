package ethereum

import (
	"fmt"
)

type Watch struct {
}

func (w *Watch) notify(peerID string, request *newBlockData) {

	a := request.Block.Number()
	b := request.TD.String()
	c := request.Block.Header.Difficulty

	fmt.Printf("===> NOTIFY (%s) Block: %d Difficulty %d. Total: %s\n", peerID, a, c, b)

}

func (w *Watch) announce(peerID string, ann []*announcement) {
	fmt.Printf("===> NOTIFY (%s) HASHES\n", peerID)
	fmt.Println(peerID, ann)
}
