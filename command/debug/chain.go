package debug

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/umbracle/minimal/blockchain"
	"github.com/umbracle/minimal/blockchain/storage/leveldb"

	"github.com/spf13/cobra"
)

var chainCmd = &cobra.Command{
	Use:   "chain",
	Short: "utility tool to debug a blockchain chain of headers",
	RunE:  chaindebugE,
}

func init() {
	DebugCmd.AddCommand(chainCmd)
}

func chaindebugE(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected one argument")
	}

	path := args[0]
	if !strings.HasSuffix(path, "blockchain") {
		path = filepath.Join(path, "blockchain")
	}

	storage, err := leveldb.NewLevelDBStorage(path, nil)
	if err != nil {
		return err
	}

	blockchain := blockchain.NewBlockchain(storage, nil, nil, nil)
	height, _ := blockchain.Header()

	for i := uint64(0); i < height.Number; i++ {
		h, ok := blockchain.GetHeaderByNumber(i)
		if !ok {
			panic("not found")
		}
		fmt.Printf("%d|%s\n", h.Number, h.Miner.String())
	}
	return nil
}
