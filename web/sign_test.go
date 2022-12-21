package main

import (
	"encoding/hex"
	"testing"
)

func TestSignGenerate(t *testing.T) {
	privKey := "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	contract := "0x5FC8d32690cc91D4c39d9d3abcBD16989F875707"
	to := "0x90F79bf6EB2c4f870365E785982E1f101E93b906"
	salt, _ := hex.DecodeString("7127573f4b02384b29684e6abaff7b307d838b62c78246eafea4931aa393dc54")

	priv, err := parsePrivateKey(privKey)
	if err != nil {
		t.Fatalf("Unable to parsePrivateKey: %v", err)
		return
	}

	signature, salt, err := signWithdraw(to, contract, "5000000000000000000", salt, priv)
	if err != nil {
		t.Fatalf("Unable to signWithdraw: %v", err)
		return
	}

	if hex.EncodeToString(signature) != "cee83f8ec41b51b15cc30e04207f8973a6b5ef82418cfe79c8b3f2843bc8d73c67be03e02bcd532cdf9409671b6ec9eb562916ebc671d7c08bd4455b55084bac1c" {
		t.Fatalf("Invalid signature: %s", hex.EncodeToString(signature))
		return
	}
}
