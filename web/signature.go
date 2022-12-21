package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"strconv"
	"strings"
)

type sendSignatureResponse struct {
	Success   bool   `json:"success"`
	Signature string `json:"signature"`
	Salt      string `json:"salt"`
	Amount    string `json:"amount"`
}

func sendSignature(ctx *controllerContext) error {
	msg := ctx.req.FormValue("msg")
	dest := ctx.req.FormValue("dest")
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
	valid, amount, err := checkTransaction(ctx.client, signer, tx)
	if err != nil {
		return fmt.Errorf("Unable to checkTransaction: %v", err)
	}
	if !valid {
		return nil
	}

	// Generate EIP-712 signature for the client
	signature, salt, err := signWithdraw(dest, "0x213C4dFfFD764765d11FbC067b9Ef89853CCb4a3", amount, nil, ctx.priv)
	if err != nil {
		return fmt.Errorf("Unable to signWithdraw: %v", err)
	}

	// Return data as json
	response, err := json.Marshal(sendSignatureResponse{
		Success:   true,
		Signature: hex.EncodeToString(signature),
		Salt:      hex.EncodeToString(salt),
		Amount:    amount,
	})
	if err != nil {
		return fmt.Errorf("Unable to Marshal JSON response: %v", err)
	}

	ctx.res.Header().Add("Content-Type", "application/json")
	ctx.res.Write(response)
	return nil
}

func checkTransaction(client *ethclient.Client, expectedSigner, txHash string) (bool, string, error) {
	tx, pending, err := client.TransactionByHash(context.Background(), common.HexToHash(txHash))
	if err != nil {
		return false, "", err
	}
	if pending {
		return false, "", nil
	}

	// Retrieve signer
	signer, err := types.LatestSignerForChainID(tx.ChainId()).Sender(tx)
	if err != nil {
		return false, "", fmt.Errorf("Unable to retrieve transaction latest signer: %v", err)
	}
	if strings.ToLower(signer.Hex()) != strings.ToLower(expectedSigner) {
		return false, "", nil
	}
	return true, tx.Value().String(), nil
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
