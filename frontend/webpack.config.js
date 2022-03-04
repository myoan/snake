const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');

// const pathToPhaser = path.join(__dirname, '/mode_modules/pahser/');
// const phaser = path.join(pathToPhaser, 'dist/phaser.js');

const {CleanWebpackPlugin} = require('clean-webpack-plugin');

module.exports = {
  mode: 'development',
  devtool: 'inline-source-map',
  entry: './app.ts',
  watch: false,
  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: 'bundle.js',
  },
  resolve: {
    // Add `.ts` and `.tsx` as a resolvable extension.
    extensions: ['.ts', '.tsx', '.js'],
  },
  devServer: {
    static: {
      directory: path.join(__dirname, 'dist')
    },
    compress: false,
    port: 8081,
  },
  plugins: [
    new CleanWebpackPlugin({verbose: true}),
    new HtmlWebpackPlugin(),
  ],
  module: {
    rules: [
      // all files with a `.ts` or `.tsx` extension will be handled by `ts-loader`
      {test: /\.tsx?$/, loader: 'ts-loader'},
    ],
  },
};
