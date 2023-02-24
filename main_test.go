package main

import (
	"syscall/js"
	"testing"

	"github.com/spacemeshos/ed25519"
	"github.com/stretchr/testify/require"
)

func Test_Generate(t *testing.T) {
	seed := make([]byte, ed25519.SeedSize)

	seedBytes := js.Global().Get("Uint8Array").New(ed25519.SeedSize)
	n := js.CopyBytesToJS(seedBytes, seed)
	require.Equal(t, ed25519.SeedSize, n)

	key := ed25519.NewKeyFromSeed(seed)

	GenerateKeyCallback.Invoke(seedBytes, js.FuncOf(func(this js.Value, args []js.Value) any {
		t.Log("GenerateKeyCallback invoked")

		pubKey := make([]byte, ed25519.PublicKeySize)
		n := js.CopyBytesToGo(pubKey, args[0])
		require.Equal(t, ed25519.PublicKeySize, n)
		require.EqualValues(t, key.Public(), pubKey)

		privKey := make([]byte, ed25519.PrivateKeySize)
		n = js.CopyBytesToGo(privKey, args[1])
		require.Equal(t, ed25519.PrivateKeySize, n)
		require.EqualValues(t, key, privKey)

		return nil
	}))
}
