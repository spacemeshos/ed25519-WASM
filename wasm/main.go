package main

import (
	"bytes"
	"syscall/js"

	"github.com/spacemeshos/ed25519"
)

var c chan bool

func init() {
	c = make(chan bool)
}

// GenerateKeyCallback is a callback function that generates a new key pair from a seed.
// Parameters:
//   1. seed - a Uint8Array of 32 bytes
//   2. callback - a function that will be called with the generated public and private keys
var GenerateKeyCallback = js.FuncOf(func(this js.Value, args []js.Value) any {
	callback := args[1]

	seed := make([]byte, ed25519.SeedSize)
	n := js.CopyBytesToGo(seed, args[0])
	if n != ed25519.SeedSize {
		callback.Invoke(js.Null(), js.Null())
		return nil
	}

	publicKey, privateKey, err := ed25519.GenerateKey(bytes.NewReader(seed))
	if err != nil {
		callback.Invoke(js.Null(), js.Null())
		return nil
	}
	pubBytes := js.Global().Get("Uint8Array").New(ed25519.PublicKeySize)
	n = js.CopyBytesToJS(pubBytes, publicKey)
	if n != ed25519.PublicKeySize {
		callback.Invoke(js.Null(), js.Null())
		return nil
	}
	privBytes := js.Global().Get("Uint8Array").New(ed25519.PrivateKeySize)
	n = js.CopyBytesToJS(privBytes, privateKey)
	if n != ed25519.PrivateKeySize {
		callback.Invoke(js.Null(), js.Null())
		return nil
	}

	callback.Invoke(pubBytes, privBytes)
	return nil
})

// DerivePrivateKeyCallback is a callback function that derives a new key pair from a seed.
// Parameters:
//   1. seed - a Uint8Array of 32 bytes
//   2. index - a number
//   3. salt - a Uint8Array of 32 bytes
//   4. callback - a function that will be called with the generated public and private keys
var DerivePrivateKeyCallback = js.FuncOf(func(this js.Value, args []js.Value) any {
	callback := args[3]

	seed := make([]byte, ed25519.SeedSize)
	n := js.CopyBytesToGo(seed, args[0])
	if n != ed25519.SeedSize {
		callback.Invoke(js.Null(), js.Null())
		return nil
	}

	index := uint64(args[1].Int())
	salt := make([]byte, args[2].Length())
	n = js.CopyBytesToGo(seed, args[2])
	if n != args[2].Length() {
		callback.Invoke(js.Null(), js.Null())
		return nil
	}

	privateKey := ed25519.NewDerivedKeyFromSeed(seed, index, salt)
	pubBytes := js.Global().Get("Uint8Array").New(ed25519.PublicKeySize)
	n = js.CopyBytesToJS(pubBytes, privateKey[32:])
	if n != ed25519.PublicKeySize {
		callback.Invoke(js.Null(), js.Null())
		return nil
	}
	privBytes := js.Global().Get("Uint8Array").New(ed25519.PrivateKeySize)
	n = js.CopyBytesToJS(privBytes, privateKey)
	if n != ed25519.PrivateKeySize {
		callback.Invoke(js.Null(), js.Null())
		return nil
	}
	callback.Invoke(pubBytes, privBytes)
	return nil
})

// Sign2Callback is a callback function that signs a message with a private key using a modified version of ed25519 signature scheme to allow key extraction.
// Parameters:
//   1. privateKey - a Uint8Array of 64 bytes
//   2. message - a Uint8Array
//   3. callback - a function that will be called with the generated signature
var Sign2Callback = js.FuncOf(func(this js.Value, args []js.Value) any {
	callback := args[2]
	
	sk := make([]byte, ed25519.PrivateKeySize)
	n := js.CopyBytesToGo(sk, args[0])
	if n != ed25519.PrivateKeySize {
		callback.Invoke(js.Null(), js.Null())
		return nil
	}
	message := make([]byte, args[1].Length())
	n = js.CopyBytesToGo(message, args[1])
	if n != args[1].Length() {
		callback.Invoke(js.Null(), js.Null())
		return nil
	}
	
	signature := ed25519.Sign2(sk, message)
	sigBytes := js.Global().Get("Uint8Array").New(ed25519.SignatureSize)
	n = js.CopyBytesToJS(sigBytes, signature)
	if n != ed25519.SignatureSize {
		callback.Invoke(js.Null(), js.Null())
		return nil
	}
	callback.Invoke(sigBytes)
	return nil
})

// Verify2Callback is a callback function that verifies a signature with a public key using a modified version of ed25519 signature scheme to allow key extraction.
// Parameters:
//   1. publicKey - a Uint8Array of 32 bytes
//   2. message - a Uint8Array
//   3. signature - a Uint8Array of 64 bytes
//   4. callback - a function that will be called with the generated signature
var Verify2Callback = js.FuncOf(func(this js.Value, args []js.Value) any {
	callback := args[3]

	pk := make([]byte, ed25519.PublicKeySize)
	n := js.CopyBytesToGo(pk, args[0])
	if n != ed25519.PublicKeySize {
		callback.Invoke(js.Null(), js.Null())
		return nil
	}
	message := make([]byte, args[1].Length())
	n = js.CopyBytesToGo(message, args[1])
	if n != args[1].Length() {
		callback.Invoke(js.Null(), js.Null())
		return nil
	}
	signature := make([]byte, ed25519.SignatureSize)
	n = js.CopyBytesToGo(signature, args[2])
	if n != ed25519.SignatureSize {
		callback.Invoke(js.Null(), js.Null())
		return nil
	}

	isValid := ed25519.Verify2(pk, message, signature)
	callback.Invoke(isValid)
	return nil
})

var ShutdownCallback = js.FuncOf(func(this js.Value, args []js.Value) any {
	c <- true
	return nil
})

func RegisterCallbacks() {
	js.Global().Set("__generateKeyPair", GenerateKeyCallback)
	js.Global().Set("__deriveNewKeyPair", DerivePrivateKeyCallback)
	js.Global().Set("__signTransaction", Sign2Callback)
	js.Global().Set("__verifyTransaction", Verify2Callback)
	js.Global().Set("__stopAndCleanUp", ShutdownCallback)
}

func CleanUp() {
	js.Global().Set("__generateKeyPair", js.Undefined())
	js.Global().Set("__deriveNewKeyPair", js.Undefined())
	js.Global().Set("__signTransaction", js.Undefined())
	js.Global().Set("__verifyTransaction", js.Undefined())
	js.Global().Set("__stopAndCleanUp", js.Undefined())
}

func main() {
	RegisterCallbacks()

	<-c

	CleanUp()
	GenerateKeyCallback.Release()
	DerivePrivateKeyCallback.Release()
	Sign2Callback.Release()
	Verify2Callback.Release()
	ShutdownCallback.Release()
}
