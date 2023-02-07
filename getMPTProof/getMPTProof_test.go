package getmptproof_test

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	getmptproof "github.com/hegelXu/GetMerkleProof/getMPTProof"
	verifymptproof "github.com/hegelXu/GetMerkleProof/verifyMPTProof"
)

// src: https://chainlist.org/chain/1
var URL = []string{
	"https://rpc.payload.de",
	"https://rpc.builder0x69.io",
	"https://ethereum.publicnode.com",
	"https://endpoints.omniatech.io/v1/eth/mainnet/public",
	"https://1rpc.io/eth",
}

func TestGetMerkleProof(t *testing.T) {
	var client *ethclient.Client
	var err error
	for _, url := range URL {
		client, err = ethclient.Dial(url)
		if err != nil {
			log.Fatal(err)
		} else {
			break
		}
	}
	defer client.Close()

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	block, err := client.BlockByNumber(context.Background(), header.Number)
	if err != nil {
		log.Fatal(err)
	}

	txIndex := rand.Intn(len(block.Transactions()))
	tx := block.Body().Transactions[txIndex]
	proof, err := getmptproof.GetMerkleProof(block, tx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("-----Test the correct proof:-----")
	if !verifymptproof.VerifyMerkleProof(tx, block.TxHash(), *proof) {
		fmt.Println("Verification failed.")
		t.Error("failed")
	} else {
		fmt.Println("Verification succeeded.")
	}

	fmt.Print("-----Test incorrect proof:-----")
	proof.Value = append(proof.Value, '!')
	if !verifymptproof.VerifyMerkleProof(tx, block.TxHash(), *proof) {
		fmt.Println("Incorrect MPT Proof.")
	} else {
		fmt.Println("wrong!")
		t.Error("failed")
	}

}
