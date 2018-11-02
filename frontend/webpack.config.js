const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin')

module.exports = (env, argv) => {

  console.log('Server Env:', env.SERVER_ENV)
  
  return ({
    context: path.resolve(__dirname, '.'),
    entry: './js/index.js',
    mode: env.SERVER_ENV === 'production' ? 'production' : 'development',
    devtool: env.SERVER_ENV === 'production' ? 'source-maps' : 'eval',
    output: {
      filename: 'static/js/main.js',
      path: path.resolve(__dirname, '../artifact/' + env.SERVER_ENV),
      publicPath: '/'
    },
    performance: {
      hints: 'error',
      maxEntrypointSize: 1200000,
      maxAssetSize: 1200000,
    },
    plugins: [
      new HtmlWebpackPlugin({
        title: 'BG Mentor',
        hash: true,
        filename: 'static/index.html'
      })
    ],
    module: {
      rules: [
        {
          test: /\.js$/,
          exclude: /node_modules/,
          use: {
            loader: 'babel-loader'
          }
        }
      ]
    }
  })
}
