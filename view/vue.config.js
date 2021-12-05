module.exports = {
    devServer: {
        port: 8989,
        open: true,
        proxy: {
            '/api': {
                target: 'http://127.0.0.1:8000/api',
                changeOrigin: true,
                ws: true,
                pathRewrite: {
                    '^/api': ''
                }
            }
        }
    }
}