import {defineConfig} from 'vite'
import reactRefresh from '@vitejs/plugin-react-refresh'

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [reactRefresh()],
    server: {
        proxy: {
            // string shorthand
            //'/things': 'http://localhost:9090',
            // with options
            '/things': {
                target: 'http://localhost:9090',
                changeOrigin: true,
                secure: false,
                ws:true,
               // rewrite: (path) => path.replace(/^\/things/, '')
            },
            '/new_things': {
                target: 'ws://localhost:9090',
                changeOrigin: false,
                secure: false,
                ws:true,
                // rewrite: (path) => path.replace(/^\/things/, '')
            },
            '/actions': {
                target: 'http://localhost:9090',
                changeOrigin: true,
                secure: false,
                ws:true,
                // rewrite: (path) => path.replace(/^\/things/, '')
            }
            // with RegEx

        }
    }
})



