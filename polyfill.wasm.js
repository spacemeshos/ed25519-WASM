(() => {
  // Map web browser API and Node.js API to a single common API (preferring web standards over Node.js API).
	const isNodeJS = global.process && /\/?node\s?/.test(global.process.title);
	if (isNodeJS) {
		global.require = require;
		global.fs = require("fs");

		const nodeCrypto = require("crypto");
		global.crypto = {
			getRandomValues(b) {
				nodeCrypto.randomFillSync(b);
			},
		};

		global.performance = {
			now() {
				const [sec, nsec] = process.hrtime();
				return sec * 1000 + nsec / 1000000;
			},
		};

		const util = require("util");
		global.TextEncoder = util.TextEncoder;
		global.TextDecoder = util.TextDecoder;
	} else {
		let outputBuf = "";
		global.fs = {
			constants: { O_WRONLY: -1, O_RDWR: -1, O_CREAT: -1, O_TRUNC: -1, O_APPEND: -1, O_EXCL: -1 }, // unused
			writeSync(fd, buf) {
				outputBuf += decoder.decode(buf);
				const nl = outputBuf.lastIndexOf("\n");
				if (nl != -1) {
					console.log(outputBuf.substr(0, nl));
					outputBuf = outputBuf.substr(nl + 1);
				}
				return buf.length;
			},
			write(fd, buf, offset, length, position, callback) {
				if (offset !== 0 || length !== buf.length || position !== null) {
					throw new Error("not implemented");
				}
				const n = this.writeSync(fd, buf);
				callback(null, n);
			},
			open(path, flags, mode, callback) {
				const err = new Error("not implemented");
				err.code = "ENOSYS";
				callback(err);
			},
			read(fd, buffer, offset, length, position, callback) {
				const err = new Error("not implemented");
				err.code = "ENOSYS";
				callback(err);
			},
			fsync(fd, callback) {
				callback(null);
			},
		};
	}
})();
