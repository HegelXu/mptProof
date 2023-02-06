package verifymptproof

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	mptproof "github.com/hegelXu/GetMerkleProof/mptProof"
)

func VerifyMerkleProof(transaction *types.Transaction, txHash common.Hash, mptp mptproof.MPTProof) bool {
	db := trie.NewDatabase(rawdb.NewMemoryDatabase())
	mpt := trie.NewEmpty(db)
	types.DeriveSha(types.Transactions{transaction}, mpt)

	var indexBuf []byte
	indexBuf = rlp.AppendUint64(indexBuf[:0], uint64(0))
	txValue := mpt.Get(indexBuf)

	val, err := trie.VerifyProof(txHash, mptp.Index, &mptp.Proof)

	if err == nil && bytes.Equal(val, mptp.Value) && bytes.Equal(val, txValue) {
		return true
	}

	return false
}
