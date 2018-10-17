module.exports = {
	devServer: {
		proxy: {
			'/img': {
				target: 'http://blackbeard.apertoire.org:7623',
			},
		},
	},
}
