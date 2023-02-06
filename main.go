package main

import (
	"fmt"
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/trie"
	getmptproof "github.com/hegelXu/GetMerkleProof/getMPTProof"
	verifymptproof "github.com/hegelXu/GetMerkleProof/verifyMPTProof"
)

func main() {
	block := newTestBlock()
	txIndex := rand.Intn(len(block.Transactions()))
	tx := block.Body().Transactions[txIndex]
	proof, err := getmptproof.GetMerkleProof(block, tx)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Print("Test the correct proof:")
	if !verifymptproof.VerifyMerkleProof(tx, block.TxHash(), *proof) {
		fmt.Println("Verification failed.")
	} else {
		fmt.Println("Verification succeeded.")
	}

	fmt.Print("Test incorrect proof:")
	proof.Value = append(proof.Value, '!')
	if !verifymptproof.VerifyMerkleProof(tx, block.TxHash(), *proof) {
		fmt.Println("Incorrect MPT Proof.")
	} else {
		fmt.Println("wrong!")
	}

}

func newTestBlock() *types.Block {
	var (
		key, _   = crypto.GenerateKey()
		txs      = make([]*types.Transaction, 70)
		receipts = make([]*types.Receipt, len(txs))
		signer   = types.LatestSigner(params.TestChainConfig)
		uncles   = make([]*types.Header, 3)
	)
	header := &types.Header{
		Difficulty: math.BigPow(11, 11),
		Number:     math.BigPow(2, 9),
		GasLimit:   12345678,
		GasUsed:    1476322,
		Time:       9876543,
		Extra:      []byte("coolest block on chain"),
	}
	for i := range txs {
		amount := math.BigPow(2, int64(i))
		price := big.NewInt(300000)
		data := make([]byte, 100)
		tx := types.NewTransaction(uint64(i), common.Address{}, amount, 123457, price, data)
		signedTx, err := types.SignTx(tx, signer, key)
		if err != nil {
			panic(err)
		}
		txs[i] = signedTx
		receipts[i] = types.NewReceipt(make([]byte, 32), false, tx.Gas())
	}
	for i := range uncles {
		uncles[i] = &types.Header{
			Difficulty: math.BigPow(11, 11),
			Number:     math.BigPow(2, 9),
			GasLimit:   12345678,
			GasUsed:    1476322,
			Time:       9876543,
			Extra:      []byte("benchmark uncle"),
		}
	}
	tr1 := trie.NewEmpty(trie.NewDatabase(rawdb.NewMemoryDatabase()))
	return types.NewBlock(header, txs, uncles, receipts, tr1)
}
