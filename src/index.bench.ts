import { benchmarkSuite } from "jest-bench";
import Ed25519 from '../src/index';

const seed = Uint8Array.from([
    1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
    1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
]);

const message = Buffer.from('Hello world!');

let pub: Uint8Array;
let sec: Uint8Array;
let signature: Uint8Array;

benchmarkSuite("Ed25519", {
    async setupSuite() {
        const k = await Ed25519.generateKeyPair(seed);
        expect(k.publicKey).toHaveLength(32);
        expect(k.secretKey).toHaveLength(64);

        pub = k.publicKey;
        sec = k.secretKey;

        signature = await Ed25519.sign(sec, message);
        expect(signature).toHaveLength(64);
    },

    ["sign"]: async () => {
        const s = await Ed25519.sign(sec, message);
        expect(s).toEqual(signature);
    },

    ["verify"]: async () => {
        const v = await Ed25519.verify(pub, message, signature);
        expect(v).toBeTruthy();
    }
});
