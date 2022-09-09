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

func TypedArrayToByteSlice(arg js.Value) []byte {
	length := arg.Length()
	bytesToReturn := make([]byte, length)
	for i := 0; i < length; i++ {
		bytesToReturn[i] = byte(arg.Index(i).Int())
	}
	return bytesToReturn
}

func SliceToJSArray(slice []byte) js.Value {
	jsArry := js.Global().Get("Uint8Array").New(len(slice))
	js.CopyBytesToJS(jsArry, slice)
	return jsArry
}

var GeneratePrivateKeyCallback = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	seed := TypedArrayToByteSlice(args[0])
	seedBuffer := bytes.NewReader(seed)
	var publicKey, privateKey, err = ed25519.GenerateKey(seedBuffer)
	if err != nil {
		return nil
	}
	_ = publicKey
	return SliceToJSArray([]byte(privateKey))
})

var GeneratePublicKeyCallback = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	seed := TypedArrayToByteSlice(args[0])
	seedBuffer := bytes.NewReader(seed)
	var publicKey, privateKey, err = ed25519.GenerateKey(seedBuffer)
	if err != nil {
		return nil
	}
	_ = privateKey
	return SliceToJSArray([]byte(publicKey))
})

var DerivePrivateKeyCallback = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	seed := TypedArrayToByteSlice(args[0])
	index := uint64(args[1].Int())
	salt := TypedArrayToByteSlice(args[2])
	var privateKey = ed25519.NewDerivedKeyFromSeed(seed, index, salt)
	publicKey := make([]byte, ed25519.PublicKeySize)
	copy(publicKey, privateKey[32:])
	return SliceToJSArray([]byte(privateKey))
})

var DerivePublicKeyCallback = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	seed := TypedArrayToByteSlice(args[0])
	index := uint64(args[1].Int())
	salt := TypedArrayToByteSlice(args[2])
	var privateKey = ed25519.NewDerivedKeyFromSeed(seed, index, salt)
	publicKey := make([]byte, ed25519.PublicKeySize)
	copy(publicKey, privateKey[32:])
	return SliceToJSArray([]byte(publicKey))
})

var Sign2Callback = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	var sk ed25519.PrivateKey = TypedArrayToByteSlice(args[0])
	var message = TypedArrayToByteSlice(args[1])
	signature := ed25519.Sign2(sk, message)
	return SliceToJSArray([]byte(signature))
})

var Verify2Callback = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	var pk ed25519.PublicKey = TypedArrayToByteSlice(args[0])
	var message = TypedArrayToByteSlice(args[1])
	var signature = TypedArrayToByteSlice(args[2])
	isValid := ed25519.Verify2(pk, message, signature)
	return isValid
})

var ShutdownCallback = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	c <- true
	return nil
})

func RegisterCallbacks() {
	js.Global().Set("__generatePrivateKey", GeneratePrivateKeyCallback)
	js.Global().Set("__generatePublicKey", GeneratePublicKeyCallback)
	js.Global().Set("__derivePrivateKey", DerivePrivateKeyCallback)
	js.Global().Set("__derivePublicKey", DerivePublicKeyCallback)
	js.Global().Set("__signTransaction", Sign2Callback)
	js.Global().Set("__verifyTransaction", Verify2Callback)
	js.Global().Set("__stopAndCleanUp", ShutdownCallback)
}

func CleanUp() {
	js.Global().Set("__generatePrivateKey", js.Undefined())
	js.Global().Set("__generatePublicKey", js.Undefined())
	js.Global().Set("__derivePrivateKey", js.Undefined())
	js.Global().Set("__derivePublicKey", js.Undefined())
	js.Global().Set("__signTransaction", js.Undefined())
	js.Global().Set("__verifyTransaction", js.Undefined())
	js.Global().Set("__stopAndCleanUp",  js.Undefined())
}

func main() {
	RegisterCallbacks()

	<-c

	CleanUp()

	GeneratePrivateKeyCallback.Release()
	GeneratePublicKeyCallback.Release()
	DerivePrivateKeyCallback.Release()
	DerivePublicKeyCallback.Release()
	Sign2Callback.Release()
	Verify2Callback.Release()
	ShutdownCallback.Release()
}
