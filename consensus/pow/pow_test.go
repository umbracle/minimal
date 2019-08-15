package pow

import (
	"context"
	"testing"

	"github.com/umbracle/minimal/types"
)

func TestPow(t *testing.T) {
	pow, _ := Factory(nil, nil)

	p := &types.Header{
		Number: 0,
	}
	h := &types.Header{
		ParentHash: p.Hash(),
		Number:     1,
	}

	pow.Prepare(p, h)
	b, err := pow.Seal(context.Background(), &types.Block{Header: h})
	if err != nil {
		t.Fatal(err)
	}

	pow.VerifyHeader(p, b.Header, false, false)
}
