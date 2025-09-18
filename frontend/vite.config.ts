import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";

export default defineConfig({
  plugins: [vue()],
  optimizeDeps: {
    include: ["@openim/client-sdk"],
  },
  server: {
    host: "0.0.0.0", // 容器内监听全部 IP
    port: 5173,
  },
});
