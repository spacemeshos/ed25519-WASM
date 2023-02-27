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

func Test_Derive(t *testing.T) {
	seed := make([]byte, ed25519.SeedSize)
	seedBytes := js.Global().Get("Uint8Array").New(ed25519.SeedSize)
	n := js.CopyBytesToJS(seedBytes, seed)
	require.Equal(t, ed25519.SeedSize, n)

	index := uint64(0)

	salt := make([]byte, 32)
	saltBytes := js.Global().Get("Uint8Array").New(len(salt))
	n = js.CopyBytesToJS(saltBytes, seed)
	require.Equal(t, len(salt), n)

	key := ed25519.NewDerivedKeyFromSeed(seed, index, salt)

	DerivePrivateKeyCallback.Invoke(seedBytes, js.ValueOf(index), saltBytes, js.FuncOf(func(this js.Value, args []js.Value) any {
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

func Test_Sign2(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)
	privBytes := js.Global().Get("Uint8Array").New(ed25519.PrivateKeySize)
	n := js.CopyBytesToJS(privBytes, priv)
	require.Equal(t, ed25519.PrivateKeySize, n)

	message := []byte("hello world")
	msgBytes := js.Global().Get("Uint8Array").New(len(message))
	n = js.CopyBytesToJS(msgBytes, message)
	require.Equal(t, len(message), n)

	Sign2Callback.Invoke(privBytes, msgBytes, js.FuncOf(func(this js.Value, args []js.Value) any {
		sig := make([]byte, ed25519.SignatureSize)
		n := js.CopyBytesToGo(sig, args[0])
		require.Equal(t, ed25519.SignatureSize, n)
		require.True(t, ed25519.Verify2(pub, message, sig))

		return nil
	}))
}

func Test_Verify2(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)
	pubBytes := js.Global().Get("Uint8Array").New(ed25519.PublicKeySize)
	n := js.CopyBytesToJS(pubBytes, pub)
	require.Equal(t, ed25519.PublicKeySize, n)

	message := []byte("hello world")
	msgBytes := js.Global().Get("Uint8Array").New(len(message))
	n = js.CopyBytesToJS(msgBytes, message)
	require.Equal(t, len(message), n)

	sig := ed25519.Sign2(priv, message)
	sigBytes := js.Global().Get("Uint8Array").New(ed25519.SignatureSize)
	n = js.CopyBytesToJS(sigBytes, sig)
	require.Equal(t, ed25519.SignatureSize, n)

	Verify2Callback.Invoke(pubBytes, msgBytes, sigBytes, js.FuncOf(func(this js.Value, args []js.Value) any {
		require.True(t, args[0].Bool())

		return nil
	}))
}
