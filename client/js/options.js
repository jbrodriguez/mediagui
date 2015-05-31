const Bacon 		= require('baconjs'),
	  R 			= require('ramda'),
	  Dispatcher 	= require('./dispatcher'),
	  api 			= require('./api')

const d = new Dispatcher()

module.exports = {
	toProperty: function(initialOptions) {
		console.log('options-before')
		return d
			.stream('setOptions')
			.scan(initialOptions, (_, newOptions) => newOptions)
			.log('options')
	},

	setOptions: function(options) {
		d.push('getOptions', options)
	}
}

