package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"strconv"
)

func sendSignature(ctx *controllerContext) error {
	msg := ctx.req.FormValue("msg")
	signer := ctx.req.FormValue("signer")
	tx := ctx.req.FormValue("tx")

	// Retrieve public key from the signed message provided
	pub, err := verifySignature(msg, tx)
	if err != nil {
		return err
	}

	// Only proceed if the public key matches
	if signer != pub {
		return nil
	}

	// Check if the transaction is correct
	checkTransaction(ctx.client, tx)
	return nil
}

func checkTransaction(client *ethclient.Client, txHash string) {
	tx, pending, err := client.TransactionByHash(context.Background(), common.HexToHash(txHash))
	fmt.Println(tx, pending, err)
	fmt.Println(tx.Value())
}

func verifySignature(signedMessage, message string) (string, error) {
	// Create hash of the message
	hashedMessage := []byte("\x19Ethereum Signed Message:\n" + strconv.Itoa(len(message)) + message)
	hash := crypto.Keccak256Hash(hashedMessage)

	decodedMessage := hexutil.MustDecode(signedMessage)

	// Handle EIp-115 not implemented
	if decodedMessage[64] == 27 || decodedMessage[64] == 28 {
		decodedMessage[64] -= 27
	}

	// Recover public key
	pub, err := crypto.SigToPub(hash.Bytes(), decodedMessage)
	if err != nil {
		return "", err
	}
	if pub == nil {
		return "", errors.New("Unable to get a public key from the message")
	}
	return crypto.PubkeyToAddress(*pub).String(), nil
}
