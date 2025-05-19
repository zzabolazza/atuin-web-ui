module.exports = {
  devServer: {
    port: 8081,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  },
  configureWebpack: {
    resolve: {
      alias: {
        '@': require('path').resolve(__dirname, 'src')
      }
    }
  },
  // Set the page title
  pages: {
    index: {
      entry: 'src/main.js',
      title: 'Atuin Web UI'
    }
  }
}
