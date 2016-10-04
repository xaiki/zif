var path = require('path')
var webpack = require('webpack')

module.exports = {
  entry: './src/main.js',

  output: {
    path: path.resolve(__dirname, './dist'),
    publicPath: '/dist/',
    filename: 'build.js'
  },

  resolveLoader: {
    root: path.join(__dirname, 'node_modules'),
  },

  module: {
    loaders: [
      {
        test: /\.vue$/,
        loader: 'vue'
      },

      {
        test: /\.js$/,
        loader: 'babel',
        exclude: /node_modules/
      },

      {
        test: /\.(png|jpg|gif|svg)$/,
        loader: 'file',
        query: {
          name: '[name].[ext]?[hash]'
        }
      }

    ]
  },

  plugins: [
        new webpack.ProvidePlugin({
           $: "jquery",
           jQuery: "jquery"
       })
    ]
}
