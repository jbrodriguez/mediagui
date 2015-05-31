const Bacon 		= require('baconjs'),
	  // R 			= require('ramda'),
	  Dispatcher 	= require('./dispatcher'),
	  api 			= require('./api')

const d = new Dispatcher()

module.exports = {
	toProperty: function(initialConfig) {
		console.log('settings-before')
		const gotConfig = d
			.stream('getConfig')
			.flatMap( name => Bacon.fromPromise( api.getConfig() ) )
			.log('settings-middle')

		return Bacon.update(
			initialConfig,
			gotConfig, (_, newConfig) => newConfig
		)
		.log('settings-final:')

	},

	// Public API
	getConfig: function() {
		console.log('getting config')
		d.push('getConfig')
	}

}