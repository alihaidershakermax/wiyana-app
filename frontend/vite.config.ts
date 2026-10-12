import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import tailwindcss from "@tailwindcss/vite";
import path from "path";

export default defineConfig({
  plugins: [vue(), tailwindcss()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
      // The published screenshot-prevention package points its `main`/`types`
      // at non-existent dist files. Resolve to the real ESM build instead.
      "screenshot-prevention": path.resolve(
        __dirname,
        "node_modules/screenshot-prevention/dist/screenshot-prevention.esm.js"
      ),
    },
  },
  server: {
    port: 3000,
    host: "0.0.0.0",
    strictPort: true,
    proxy: {
      "/api": {
        target: "http://localhost:8000",
        changeOrigin: true,
        secure: false,
      },
    },
  },
});
