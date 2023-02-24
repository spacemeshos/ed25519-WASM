package main

import (
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
		t.Log("args[0]:", args[0])
		t.Log("args[1]:", args[1])
		return nil
	}))
}
