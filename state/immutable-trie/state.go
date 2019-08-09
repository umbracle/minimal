package itrie

import (
	"fmt"

	lru "github.com/hashicorp/golang-lru"
	"github.com/umbracle/minimal/state"
	"github.com/umbracle/minimal/types"
)

type State struct {
	storage Storage
	cache   *lru.TwoQueueCache
}

func NewState(storage Storage) *State {

	// cache, _ := lru.NewARC(128)
	cache, _ := lru.New2Q(128)

	s := &State{
		storage: storage,
		cache:   cache,
	}
	return s
}

func (s *State) Do() {

	// s.NewSnapshotAt(types.StringToHash("0x45034d234555cff5a436be68f3096dcc012197623c9532c6bbeea24617ffe619"))
}

func (s *State) NewSnapshot() state.Snapshot {
	t := NewTrie()
	t.state = s
	t.storage = s.storage
	return t
}

func (s *State) SetCode(hash types.Hash, code []byte) {
	s.storage.SetCode(hash, code)
}

func (s *State) GetCode(hash types.Hash) ([]byte, bool) {
	return s.storage.GetCode(hash)
}

func (s *State) NewSnapshotAt(root types.Hash) (state.Snapshot, error) {

	/*
		correct := false
		if root.String() == "0x45034d234555cff5a436be68f3096dcc012197623c9532c6bbeea24617ffe619" {
			correct = true
			fmt.Printf("New snapshot at: %s\n", root.String())
		}
	*/

	tt, ok := s.cache.Get(root)
	if ok {
		t := tt.(*Trie)
		t.state = s

		/*
			if correct {
				fmt.Println(&(tt.(*Trie).root))
				fmt.Println(tt.(*Trie).root.(*FullNode))
			}
		*/

		return tt.(*Trie), nil
	}
	n, ok, err := GetNode(root.Bytes(), s.storage)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	t := &Trie{
		root:    n,
		state:   s,
		storage: s.storage,
	}
	return t, nil
}

func (s *State) AddState(root types.Hash, t *Trie) {
	/*
		if root.String() == "0x45034d234555cff5a436be68f3096dcc012197623c9532c6bbeea24617ffe619" {
			fmt.Println("==> ADD STATE")
			fmt.Println(&t.root)
			fmt.Println(t.root.(*FullNode))
		}
	*/
	s.cache.Add(root, t)
}
