import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import * as path from 'node:path'

export default defineConfig({
    plugins: [react()],
    resolve: {
        alias: {
            '@': path.resolve(__dirname, 'src'),
        },
    },
    server: {
        allowedHosts: true,
        proxy: {
            '/api': {
                target: 'http://localhost:8081', // backend
                changeOrigin: true,
                rewrite: path => path.replace(/^\/api/, ''),  // убираем "/api" перед запросом

            },
        },
    },
})
