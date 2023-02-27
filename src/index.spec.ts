import Ed25519 from '../src/index';

describe('@spacemesh/ed25519', () => {
  const seed = Uint8Array.from([
    1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
    1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
  ]);
  const salt = seed.map((x) => Math.floor((x * 2) ** 2 / 3) + 3);

  it('generateKeyPair', async () => {
    const k = await Ed25519.generateKeyPair(seed);
    expect(k.publicKey).toHaveLength(32);
    expect(k.secretKey).toHaveLength(64);
  });

  it('deriveNewKeyPair', async () => {
    const k = await Ed25519.generateKeyPair(seed);
    const k0 = await Ed25519.deriveNewKeyPair(seed, 0, salt);
    expect(k0.publicKey).toHaveLength(32);
    expect(k0.secretKey).toHaveLength(64);
    expect(k0.secretKey).not.toEqual(k.secretKey);
  });

  it('sign & verify', async () => {
    const k = await Ed25519.generateKeyPair(seed);
    const k0 = await Ed25519.deriveNewKeyPair(seed, 0, salt);
    const message = Buffer.from('Hello world!');

    const s = await Ed25519.sign(k.secretKey, message);
    expect(s).toHaveLength(64);

    const v = await Ed25519.verify(k.publicKey, message, s);
    expect(v).toBeTruthy();

    const iv = await Ed25519.verify(k0.publicKey, message, s);
    expect(iv).toBeFalsy();
  });
});
