// vite.config.ts
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      "/admin": {
        target: "http://localhost:8083",
        changeOrigin: true,
        secure: false,
        rewrite: (path) => path.replace(/^\/admin/, ""),
      },
    },
  },
});
