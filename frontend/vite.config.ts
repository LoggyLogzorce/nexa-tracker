import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    allowedHosts: ['task.nexa.test'],
    proxy: {
      '/api': {
        target: 'http://192.168.3.69:8080',
        changeOrigin: true,
      },
      '/uploads': {
        target: 'http://192.168.3.69:8080',
        changeOrigin: true,
      },
    },
  },
})
