// vite.config.ts
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      // Проксируем все запросы, начинающиеся с /orders, на сервис заказов
      // Любой запрос вида http://localhost:5173/orders/... → http://localhost:8082/orders/...
      "/orders": {
        target: "http://localhost:8082",
        changeOrigin: true,
        secure: false,
        // Без переписывания пути, т.к. API на сервисе слушает /orders:
        // rewrite: (path) => path.replace(/^\/orders/, "/orders"),
      },

      // Проксируем все запросы, начинающиеся с /products, на сервис товаров
      // http://localhost:5173/products/... → http://localhost:8081/products/...
      "/products": {
        target: "http://localhost:8081",
        changeOrigin: true,
        secure: false,
        // rewrite: (path) => path.replace(/^\/products/, "/products"),
      },

      // Проксируем все запросы, начинающиеся с /admin, на сервис админки
      // Например: http://localhost:5173/admin/login → http://localhost:8083/login
      "/admin": {
        target: "http://localhost:8083",
        changeOrigin: true,
        secure: false,
        rewrite: (path) => path.replace(/^\/admin/, ""),
      },
    },
  },
});
