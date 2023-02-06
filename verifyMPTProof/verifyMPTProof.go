package verifymptproof

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/trie"
	mptproof "github.com/hegelXu/GetMerkleProof/mptProof"
)

func VerifyMerkleProof(txHash common.Hash, mptp mptproof.MPTProof) bool {
	val, err := trie.VerifyProof(txHash, mptp.Index, &mptp.Proof)

	if err == nil && bytes.Equal(val, mptp.Value) {
		return true
	}

	return false
}
