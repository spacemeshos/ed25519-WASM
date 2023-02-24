//
// This is auto-generated file.
// Run `./scripts/inlinewasm.js some.wasm` to re-generate it
//
/* eslint-disable prettier/prettier */

const encoded = '${base64data}';
const decoded = Buffer.from(encoded, 'base64');
// const len = decoded.length;
// const bytes = new Uint8Array(len);
// console.log('?', decoded);

// for (var i = 0; i < len; i++) {
//     bytes[i] = decoded.charCodeAt(i);
// }

module.exports = () => decoded;

// const encoded = '${base64data}';

// module.exports = () =>
//   Promise.resolve(
//     new Response(
//       Buffer.from(encoded, 'base64'),
//       {
//         status: 200,
//         headers: {
//           'Content-Type': 'application/wasm',
//         },
//       }
//     )
//   );
