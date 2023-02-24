// Assume add.wasm file exists that contains a single function adding 2 provided arguments
import { readFileSync } from 'fs';

const wasmBuffer = readFileSync('ed25519.wasm');
WebAssembly.instantiate(wasmBuffer).then(wasmModule => {
    // Exported function live under instance.exports
    const { __generateKeyPair } = wasmModule.instance.exports;
    const { __deriveNewKeyPair } = wasmModule.instance.exports;
    const { __signTransaction } = wasmModule.instance.exports;
    const { __verifyTransaction } = wasmModule.instance.exports;
    const { __stopAndCleanUp } = wasmModule.instance.exports;

    const ret = __generateKeyPair(this, {}, (publicKey, privateKey) => { console.log(publicKey, privateKey) });
    console.log(ret);
});