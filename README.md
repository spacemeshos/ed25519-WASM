# ed25519-WASM

A drop-in replacement for `golang/crypto/ed25519` ([godoc](https://godoc.org/golang.org/x/crypto/ed25519),
[github](https://github.com/golang/crypto/tree/master/ed25519))
 with additional functionality compiled to WASM

## Motivation

In order to verify the validity of a given signature, the validator should posses the public key of the signer. It can be sent along with the message and its signature, which means that the overall data being sent includes 256 bits of the public key. Our function allows to extract the public key from the signature (and the message), thus the public key may not be sent, resulting in a smaller transferred data. Note: there's a computational cost for extracting the public key, so one should consider the trade-off between computations and data size.

## Usage

```bash
npm install "github.com/spacemeshos/ed25519-wasm"
```

Internal usage of package `ed25519` from `github.com/spacemeshos/ed25519` instead of `crypto/ed25519`.

## __generateKeyPair

__generateKeyPair is JS exposure of golang func GenerateKey which returns privateKey and publicKey using provided seed.

```bash
__generateKeyPair(seed Uint8Array(32), callbackToStoreValues Function) publicKey Uint8Array(32), privateKey Uint8Array(64)
```

## __signTransaction

__signTransaction is JS exposure of golang func Sign2 signs the message with privateKey and returns a signature.
The signature may be verified using Verify2(), if the signer's public key is known.
The signature returned by this method can be used together with the message
to extract the public key using ExtractPublicKey()
It will return null if privateKey.length is not PrivateKeySize.

```bash
__generateKeyPair
__signTransaction(privateKey Uint8Array(64), message Uint8Array, callbackToStoreValues Function) Uint8Array(64)
```

## __verifyTransaction

__verifyTransaction is JS exposure of Verify2 verifies a signature created with Sign2(), assuming the verifier possesses the public key.

```bash
func Verify2(publicKey PublicKey, message, sig []byte) bool
__verifyTransaction(publicKey Uint8Array(32), message Uint8Array, signature Uint8Array(64), callbackToStoreValue Function) boolean
```

## Compilation

First time usage - copy the go-js environment file to your working directory

```bash
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
```

Then compile go to WASM

```bash
GOOS=js GOARCH=wasm go build -o ed25519.wasm
```

## Go version used for compilation 1.19.6
