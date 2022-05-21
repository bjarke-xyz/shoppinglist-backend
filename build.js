const alias = require("esbuild-plugin-alias");
const {
  NodeModulesPolyfillPlugin,
} = require("@esbuild-plugins/node-modules-polyfill");

require("esbuild")
  .build({
    entryPoints: ["src/index.ts"],
    format: "iife",
    bundle: true,
    outfile: "dist/index.js",
    inject: ["./process-shim.js"],
    define: {
      // @ts-ignore
      __dirname: JSON.stringify(__dirname),
    },
    plugins: [
      NodeModulesPolyfillPlugin(),
      // @ts-ignore
      alias({
        // @ts-ignore
        "@prisma/client": require.resolve("@prisma/client"),
      }),
    ],
  })
  // @ts-ignore
  .catch(() => process.exit(1));
