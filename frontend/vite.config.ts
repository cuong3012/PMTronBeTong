import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { fileURLToPath, URL } from 'node:url'

// Cấu hình Vite cho frontend React.
// Dev server chạy cổng 5173 (khớp CORS_ORIGIN mặc định của backend).
// Proxy /api sang backend Go để tránh phải cấu hình CORS khi phát triển.
export default defineConfig({
  plugins: [react()],
  resolve: {
    // Alias '@' trỏ tới src/, đồng bộ với tsconfig.json.
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    port: 5173,
    proxy: {
      // Proxy API nghiệp vụ sang backend Go.
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      // Proxy endpoint kiểm tra sức khỏe backend (nằm ở gốc, không dưới /api).
      '/health': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
