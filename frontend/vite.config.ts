import {defineConfig} from 'vite';
import react from '@vitejs/plugin-react';

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [react()],
    build: {
        outDir: '../web',
        emptyOutDir: true,
    },
    server: {
        proxy: {
            '/api': 'http://localhost:8080',
        },
    },
});
