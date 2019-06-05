package types

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	newRlp "github.com/umbracle/go-rlp"
	"github.com/umbracle/minimal/helper/hex"
)

func TestHeaders(t *testing.T) {
	datax := "0xf901fca00000000000000000000000000000000000000000000000000000000000000000a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347948888f1f195afa192cfee860698584c030f4c9db1a0767083c42d099d13254c2287c98e48b8582c223795f44e04fd93083b6729d3b8a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421b90100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008302000080887fffffffffffffff808454c98c8142a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421880102030405060708"

	var obj Header
	if err := rlp.DecodeBytes(hex.MustDecodeHex(datax), &obj); err != nil {
		panic(err)
	}

	fmt.Println("-- header --")
	fmt.Println(obj)

	data, err := rlp.EncodeToBytes(obj)
	if err != nil {
		panic(err)
	}

	buf2, _ := newRlp.EncodeToBytes(obj)
	if !bytes.Equal(data, buf2) {

		fmt.Println(data)
		fmt.Println(buf2)

		fmt.Println("-- obj --")
		fmt.Println(obj)

		fmt.Println(hex.EncodeToString(data))

		panic("xx")
	}

	fmt.Println(data)

}

func TestEIP55(t *testing.T) {
	t.Skip("TODO")

	cases := []struct {
		address  string
		expected string
	}{
		{
			"0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed",
			"0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed",
		},
		{
			"0xfb6916095ca1df60bb79ce92ce3ea74c37c5d359",
			"0xfB6916095ca1df60bB79Ce92cE3Ea74c37c5d359",
		},
		{
			"0xdbf03b407c01e7cd3cbea99509d93f8dddc8c6fb",
			"0xdbF03B407c01E7cD3CBea99509d93f8DDDC8C6FB",
		},
		{
			"0xd1220a0cf47c7b9be7a2e6ba89f429762e7b9adb",
			"0xD1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDb",
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			addr := StringToAddress(c.address)
			assert.Equal(t, addr.String(), c.expected)
		})
	}
}
