package main

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"strings"
)

type privateKey struct {
	Private *ecdsa.PrivateKey
	Public  *ecdsa.PublicKey
	From    common.Address
}

func parsePrivateKey(priv string) (*privateKey, error) {
	privKey, err := crypto.HexToECDSA(strings.TrimPrefix(priv, "0x"))
	if err != nil {
		return nil, fmt.Errorf("Unable to parse key to ECDSA: %v", err)
	}

	// Get public key
	pubKey := privKey.Public()
	pubKeyECDSA, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("Unable to convert public key to ECDSA")
	}

	// Get FROM
	from := crypto.PubkeyToAddress(*pubKeyECDSA)

	return &privateKey{
		Private: privKey,
		Public:  pubKeyECDSA,
		From:    from,
	}, nil
}
