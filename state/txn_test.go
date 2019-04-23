package state

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	iradix "github.com/hashicorp/go-immutable-radix"
)

type mockState struct {
	snapshots map[common.Hash]TrieX
}

func (m *mockState) NewTrieAt(root common.Hash) (TrieX, error) {
	t, ok := m.snapshots[root]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return t, nil
}

func (m *mockState) NewTrie() TrieX {
	return &mockSnapshot{data: map[string][]byte{}}
}

func (m *mockState) GetCode(hash common.Hash) ([]byte, bool) {
	panic("Not implemented in tests")
}

type mockSnapshot struct {
	data map[string][]byte
}

func (m *mockSnapshot) Get(k []byte) ([]byte, bool) {
	v, ok := m.data[hexutil.Encode(k)]
	return v, ok
}

func (m *mockSnapshot) Commit(x *iradix.Tree) (TrieX, []byte) {
	panic("Not implemented in tests")
}

func newStateWithPreState(preState map[common.Address]*PreState) (*mockState, *mockSnapshot) {
	state := &mockState{
		snapshots: map[common.Hash]TrieX{},
	}
	snapshot := &mockSnapshot{
		data: map[string][]byte{},
	}
	for addr, p := range preState {
		account, snap := buildMockPreState(p)
		if snap != nil {
			state.snapshots[account.Root] = snap
		}

		accountRlp, err := rlp.EncodeToBytes(account)
		if err != nil {
			panic(err)
		}
		snapshot.data[addr.String()] = accountRlp
	}

	return state, snapshot
}

func newTestTxn(p map[common.Address]*PreState) *Txn {
	state, snap := newStateWithPreState(p)

	auxState := &State{
		state: state,
	}
	auxSnap := &Snapshot{
		state: auxState,
		tt:    snap,
	}
	return newTxn(auxState, auxSnap)
}

func buildMockPreState(p *PreState) (*Account, *mockSnapshot) {
	var snap *mockSnapshot
	root := emptyStateHash

	if p.State != nil {
		data := map[string][]byte{}
		for k, v := range p.State {
			vv, _ := rlp.EncodeToBytes(bytes.TrimLeft(v.Bytes(), "\x00"))
			data[k.String()] = vv
		}
		root = randomHash()
		snap = &mockSnapshot{
			data: data,
		}
	}

	account := &Account{
		Nonce:   p.Nonce,
		Balance: big.NewInt(int64(p.Balance)),
		Root:    root,
	}
	return account, snap
}

const letterBytes = "0123456789ABCDEF"

func randomHash() common.Hash {
	b := make([]byte, common.HashLength)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return common.BytesToHash(b)
}

func TestTxn(t *testing.T) {
	txn := newTestTxn(defaultPreState)

	txn.SetState(addr1, hash1, hash1)
	if txn.GetState(addr1, hash1) != hash1 {
		t.Fail()
	}

	ss := txn.Snapshot()
	txn.SetState(addr1, hash1, hash2)
	if txn.GetState(addr1, hash1) != hash2 {
		t.Fail()
	}

	txn.RevertToSnapshot(ss)
	if txn.GetState(addr1, hash1) != hash1 {
		t.Fail()
	}
}
