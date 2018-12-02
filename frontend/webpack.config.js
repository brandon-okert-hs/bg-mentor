const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin')

module.exports = (env, argv) => {

  console.log('Server Env:', env.SERVER_ENV)
  
  return ({
    context: path.resolve(__dirname, '.'),
    entry: {
      index: './js/index.js',
      tournaments: './js/tournaments.js',
      tournament: './js/tournament.js'
    },
    mode: env.SERVER_ENV === 'production' ? 'production' : 'development',
    devtool: env.SERVER_ENV === 'production' ? 'source-maps' : 'eval',
    output: {
      filename: 'static/js/[name].js',
      path: path.resolve(__dirname, '../artifact/' + env.SERVER_ENV),
      publicPath: '/'
    },
    performance: {
      hints: 'error',
      maxEntrypointSize: 1800000,
      maxAssetSize: 1800000,
    },
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
