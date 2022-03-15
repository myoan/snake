const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const {CleanWebpackPlugin} = require('clean-webpack-plugin');

module.exports = {
  entry: './app.ts',
  watch: false,
  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: 'bundle.js',
  },
  resolve: {
    extensions: ['.ts', '.tsx', '.js'],
  },
  plugins: [
    new CleanWebpackPlugin({verbose: true}),
    new HtmlWebpackPlugin(),
  ],
  module: {
    rules: [
      {test: /\.tsx?$/, loader: 'ts-loader'},
    ],
  },
};
