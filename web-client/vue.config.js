module.exports = {
    outputDir: "../public/admin",
    publicPath: process.env.NODE_ENV === 'production' ? '/admin/' : '/',
    devServer: {
        proxy: 'http://localhost:8080'
    }
}