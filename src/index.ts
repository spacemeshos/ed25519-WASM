type GlobalWithGoWasm = typeof globalThis & {
  Go: {
    new(): {
      run: (instance: WebAssembly.Instance) => void;
      importObject: WebAssembly.Imports;
    }
  }
};

const load = (): GlobalWithGoWasm => {
  require('../polyfill.wasm.js');
  require('../build/wasm_exec.js');

  return global as GlobalWithGoWasm;
};

const runWasm = (wasmBytes: () => Uint8Array) => {
  const bytes = wasmBytes() as Uint8Array;
  const globalGo = load();
  const go = new globalGo.Go();
  const mod = new WebAssembly.Module(bytes.buffer);
  const instance = new WebAssembly.Instance(mod, go.importObject);
  go.run(instance);
  return instance;
};

export type KeyPair = {
  publicKey: Uint8Array; // 32 bytes
  secretKey: Uint8Array; // 64 bytes
}

class Ed25519 {
  private static instance: WebAssembly.Instance | null = null;

  static init() {
    if (!this.instance) {
      this.instance = runWasm(require('../build/main.inl.js'));
    }
  }
  static cleanup() {
    if (this.instance) {
      global.__stopAndCleanUp();
      this.instance = null;
    }
  }

  static async generateKeyPair(seed: Uint8Array): Promise<KeyPair> {
    assertByteLength(32, seed);
    this.init();
    return new Promise((resolve, reject) => {
      handleCallTimeout(1000, reject);
      global.__generateKeyPair(
        seed,
        ofKeyPair(resolve, reject)
      );
    });
  }

  static deriveNewKeyPair(
    seed: Uint8Array, // 32 bytes
    index: number,
    salt: Uint8Array,  // 32 bytes
  ): Promise<KeyPair> {
    assertByteLength(32, seed);
    assertByteLength(32, salt);
    this.init();
    return new Promise<KeyPair>((resolve, reject) => {
      global.__deriveNewKeyPair(
        seed, index, salt,
        ofKeyPair(resolve, reject)
      );
    });
  }
  
  static sign(
    privateKey: Uint8Array, // 64 bytes
    message: Uint8Array,
  ): Promise<Uint8Array> {
    assertByteLength(64, privateKey);
    this.init();
    return new Promise<Uint8Array>((resolve, reject) => {
      global.__signTransaction(
        privateKey, message,
        ofBytes(resolve, reject)
      );
    });
  }
  
  static verify(
    publicKey: Uint8Array, // 32 bytes
    message: Uint8Array,
    signature: Uint8Array, // 64 bytes
  ): Promise<boolean> {
    assertByteLength(32, publicKey);
    assertByteLength(64, signature);
    this.init();
    return new Promise<boolean>((resolve, reject) => {
      global.__verifyTransaction(
        publicKey, message, signature,
        ofBoolean(resolve, reject)
      );
    });
  }
}

function handleCallTimeout(ms: number, reject: (reason: any) => void) {
  return setTimeout(
    () => reject(new Error('Can not generate key: reached timeout')),
    ms
  );
}

function assertByteLength(size: number, bytes: Uint8Array) {
  if (bytes.length < size) {
    throw new Error(
      `Expected to have ${size} bytes length, got ${bytes.length}: ${bytes}`
    );
  }
}

function ofKeyPair(resolve, reject) {
  return (publicKey, secretKey) => {
    if (!!publicKey && !!secretKey) {
      resolve({ publicKey, secretKey })
    } else {
      reject(new Error('Unknown error'))
    }
  };
}

function ofBytes(resolve, reject) {
  return (bytes) => {
    if (!!bytes && bytes instanceof Uint8Array) {
      resolve(bytes);
    } else {
      reject(new Error('Unknown error'));
    }
  }
}

function ofBoolean(resolve, reject) {
  return (x) => {
    if (x !== null) {
      resolve(Boolean(x));
    } else {
      reject(new Error('Unknown error'));
    }
  }
}

export default Ed25519;