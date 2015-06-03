const Bacon 		= require('baconjs'),
	  R 			= require('ramda'),
	  Dispatcher 	= require('./dispatcher'),
	  api 			= require('./api'),
	  storage		= require('./storage')

const d = new Dispatcher()

module.exports = {
	toProperty: function(initialOptions) {
		console.log('options-before', initialOptions)
		// const gotOptions = 
		// 	.scan(initialOptions, (_, newOptions) => newOptions)
		// 	.log('options')

		return Bacon.update(
			initialOptions,
			[d.stream('setOptions')], doSetOptions
		)

		function doSetOptions(options) {
			return options
		}
	},

	// PUBLIC API
	getOptions: function() {
		d.push('setOptions')
	}
}

