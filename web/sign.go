package main

import (
	"crypto/rand"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"math/big"
)

// Generate a valid EIP712 signed message
func signWithdraw(to, contract, amount string, salt []byte, priv *privateKey) ([]byte, []byte, error) {
	typedData, salt, err := hashTypedData(salt, amount, to, contract)
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to hashTypedData: %v", err)
	}

	// Generate signature
	signature, err := crypto.Sign(typedData, priv.Private)
	if err != nil {
		return nil, nil, err
	}

	// Legacy V
	signature[64] += 27

	return signature, salt, err
}

func hashTypedData(salt []byte, amount, to, contract string) ([]byte, []byte, error) {
	// Generate salt
	if salt == nil {
		s, err := secureRandomString(64)
		if err != nil {
			return nil, nil, fmt.Errorf("Unable to secureRandomString: %v", err)
		}
		salt = crypto.Keccak256(s)
	}

	data := apitypes.TypedData{
		Types: apitypes.Types{
			"WithdrawAction": []apitypes.Type{
				{
					Name: "amount", Type: "uint256",
				},
				{
					Name: "salt", Type: "bytes32",
				},
				{
					Name: "to", Type: "address",
				},
			},
			"EIP712Domain": []apitypes.Type{
				{
					Name: "name", Type: "string",
				},
				{
					Name: "verifyingContract", Type: "address",
				},
			},
		},
		PrimaryType: "WithdrawAction",
		Domain: apitypes.TypedDataDomain{
			Name:              "SimpleMixer",
			VerifyingContract: contract,
		},
		Message: apitypes.TypedDataMessage{
			"amount": amount,
			"salt":   salt,
			"to":     to,
		},
	}

	// Hash value
	hash, _, err := apitypes.TypedDataAndHash(data)
	return hash, salt, err
}

// Generates a random string with data from a secure PRNG
func secureRandomString(n int) ([]byte, error) {
	v := "abcdefghjiklmnopqrstuvwxymz0123456789"
	r := make([]byte, n)

	for i := 0; i < n; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(v))))
		if err != nil {
			return nil, err
		}
		r[i] = v[index.Int64()]
	}

	return r, nil
}
