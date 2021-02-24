const {createProxyMiddleware} = require('http-proxy-middleware');

module.exports = function (app) {
    app.use(
        '/addons',
        createProxyMiddleware({
            target: 'http://localhost:9090',
            changeOrigin: true,
        })
    );
    app.use(
        '/settings/addons_info',
        createProxyMiddleware({
            target: 'http://localhost:9090',
            changeOrigin: true,
        })
    );
    app.use(
        '/actions',
        createProxyMiddleware({
            target: 'http://localhost:9090',
            changeOrigin: true,
        })
    );
    app.use(
        '/things',
        createProxyMiddleware({
            target: 'http://localhost:9090',
            changeOrigin: true,
        })
    );
    app.use(
        '/things',
        createProxyMiddleware({
            target: 'ws://localhost:9090',
            changeOrigin: true,
        })
    );
    app.use(
        '/new_things',
        createProxyMiddleware({
            target: 'ws://localhost:9090',
            changeOrigin: true,
        })
    );
};
