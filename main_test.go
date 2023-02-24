package main

import (
	"encoding/hex"
	"syscall/js"
	"testing"

	"github.com/spacemeshos/ed25519"
	"github.com/stretchr/testify/require"
)

func Test_Generate(t *testing.T) {
	seed := js.Global().Get("Uint8Array").New(ed25519.SeedSize)
	n := js.CopyBytesToJS(seed, make([]byte, ed25519.SeedSize))
	require.Equal(t, ed25519.SeedSize, n)

	GenerateKeyCallback.Invoke(seed, js.FuncOf(func(this js.Value, args []js.Value) any {
		t.Log("GenerateKeyCallback invoked")

		pubKey := make([]byte, ed25519.PublicKeySize)
		n := js.CopyBytesToGo(pubKey, args[0])
		require.Equal(t, ed25519.PublicKeySize, n)
		t.Log("args[0]:", hex.EncodeToString(pubKey))

		privKey := make([]byte, ed25519.PrivateKeySize)
		n = js.CopyBytesToGo(privKey, args[1])
		require.Equal(t, ed25519.PrivateKeySize, n)
		t.Log("args[1]:", hex.EncodeToString(privKey))

		return nil
	}))
}
