var webpack = require('webpack');
var path = require('path');

var config = {
	entry: [
		'./src/index'
	],
	module: {
		loaders: [
			{ test: /\.js?$/, loader: 'babel', exclude: /node_modules/ },
			{ test: /\.s?css$/, loader: 'style!css!sass' },
			{ test: /\.json$/, loader: 'json-loader' }
		]
	},
	node: {
		fs: "empty"
	},
	resolve: {
		extensions: ['', '.js']
	},
	output: {
		path: path.join(__dirname, '/dist'),
		publicPath: '/',
		filename: 'bundle.js'
	},
	devServer: {
		contentBase: './dist',
		hot: true
	},
	plugins: [
		new webpack.optimize.OccurenceOrderPlugin(),
		new webpack.HotModuleReplacementPlugin(),
		new webpack.NoErrorsPlugin()
	]
};

config.target = webpackTargetElectronRenderer(config);

module.exports = config;
