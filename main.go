package main

import (
	"bytes"
	"github.com/spacemeshos/ed25519"
	"syscall/js"
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

var GenerateKeyCallback = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	seed := TypedArrayToByteSlice(args[0])
	seedBuffer := bytes.NewReader(seed)
	var publicKey, privateKey, err = ed25519.GenerateKey(seedBuffer)
	callback := args[1]
	if err != nil {
		callback.Invoke(js.Null(), js.Null())
		return nil
	}
	callback.Invoke(js.TypedArrayOf([]byte(publicKey)), js.TypedArrayOf([]byte(privateKey)))
	return nil
})

var Sign2Callback = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	var sk ed25519.PrivateKey = TypedArrayToByteSlice(args[0])
	var message = TypedArrayToByteSlice(args[1])
	callback := args[len(args)-1:][0]
	signature := ed25519.Sign2(sk, message)
	callback.Invoke(js.TypedArrayOf(signature))
	return nil
})

var Verify2Callback = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	var pk ed25519.PublicKey = TypedArrayToByteSlice(args[0])
	var message = TypedArrayToByteSlice(args[1])
	var signature = TypedArrayToByteSlice(args[2])
	isValid := ed25519.Verify2(pk, message, signature)
	callback := args[len(args)-1:][0]
	callback.Invoke(isValid)
	return nil
})

var ShutdownCallback = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	c <- true
	return nil
})

func RegisterCallbacks() {
	js.Global().Set("__generateKeyPair", GenerateKeyCallback)
	js.Global().Set("__signTransaction", Sign2Callback)
	js.Global().Set("__verifyTransaction", Verify2Callback)
	js.Global().Set("__stopAndCleanUp", ShutdownCallback)
}

func CleanUp() {
	js.Global().Set("__generateKeyPair", js.Undefined())
	js.Global().Set("__signTransaction", js.Undefined())
	js.Global().Set("__verifyTransaction", js.Undefined())
	js.Global().Set("__stopAndCleanUp", js.Undefined())
}

func main() {
	RegisterCallbacks()

	<-c

	CleanUp()
	GenerateKeyCallback.Release()
	Sign2Callback.Release()
	Verify2Callback.Release()
	ShutdownCallback.Release()
}
