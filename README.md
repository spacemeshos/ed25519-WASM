# This library has been deprecated and is no longer used nor maintained. Please use a standard ed25519 library instead.

# ed25519-WASM

A WASM implementation of the modified ed25519 signature algorithm with key extraction
([github](https://github.com/spacemeshos/ed25519-recovery)).

## Usage

```bash
npm install "github.com/spacemeshos/ed25519-wasm"
```

Internally this package uses `github.com/spacemeshos/ed25519-recovery`.

## __generateKeyPair

`__generateKeyPair` is the JS exported Go function `GenerateKey` which returns a privateKey and publicKey using provided seed.

```bash
__generateKeyPair(seed Uint8Array(32), callbackToStoreValues Function) publicKey Uint8Array(32), privateKey Uint8Array(64)
```

## __signTransaction

`__signTransaction` is the JS exported Go function `Sign` to sign a message with privateKey and return a signature.
The signature may be verified using `__verifyTransaction`, if the signer's public key is known.

```bash
__signTransaction(privateKey Uint8Array(64), message Uint8Array, callbackToStoreValues Function) Uint8Array(64)
```

## __verifyTransaction

`__verifyTransaction` is the JS exported Go function `Verify` to verify a signature created with `__signTransaction`,
assuming the verifier possesses the public key.

```bash
__verifyTransaction(publicKey Uint8Array(32), message Uint8Array, signature Uint8Array(64), callbackToStoreValue Function) boolean
```

## Building from sources

- Go version used for compilation 1.19.6
- NodeJS 16.15.0+
- yarn

Run

```bash
make build-all
```
