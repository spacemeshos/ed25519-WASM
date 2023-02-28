/** @type {import('ts-jest/dist/types').InitialOptionsTsJest} */
module.exports = {
  preset: 'ts-jest',
  testEnvironment: "jest-bench/environment",
  testEnvironmentOptions: {
    testEnvironment: "jest-environment-node",
    testEnvironmentOptions: {}
  },
  reporters: ["default", "jest-bench/reporter"],
  testRegex: "(/__benchmarks__/.*|\\.bench)\\.(ts|tsx|js)$",
};