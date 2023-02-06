package getmptproof

import (
	"errors"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	mptproof "github.com/hegelXu/GetMerkleProof/mptProof"
)

func GetMerkleProof(block *types.Block, transaction *types.Transaction) (*mptproof.MPTProof, error) {
	txs := block.Body().Transactions
	txTarget := transaction.Hash()
	txIndex := -1

	for i, tx := range txs {
		if tx.Hash() == txTarget {
			txIndex = i
		}
	}

	if txIndex == -1 {
		return nil, errors.New("transaction does not exist")
	}

	db := trie.NewDatabase(rawdb.NewMemoryDatabase())
	mpt := trie.NewEmpty(db)
	types.DeriveSha(block.Transactions(), mpt)

	var indexBuf []byte
	indexBuf = rlp.AppendUint64(indexBuf[:0], uint64(txIndex))
	valueTarget := mpt.Get(indexBuf)

	proof := mptproof.New()
	err := mpt.Prove(indexBuf, 0, &proof)
	if err != nil {
		return nil, err
	}

	res := mptproof.MPTProof{
		Index: indexBuf,
		Value: valueTarget,
		Proof: proof,
	}

	return &res, nil
}
