//
// This is auto-generated file.
// Run `./scripts/inlinewasm.js some.wasm` to re-generate it
//
/* eslint-disable prettier/prettier */

const encoded = '${base64data}';
const decoded = Buffer.from(encoded, 'base64');

module.exports = () => decoded;

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
