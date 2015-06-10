const Bacon 		= require('baconjs'),
	  R 			= require('ramda'),
	  Dispatcher 	= require('./dispatcher'),
	  api 			= require('./api'),
	  storage		= require('./storage')

const d = new Dispatcher()

module.exports = {
	// PUBLIC API
	setSortBy: function(sortBy) {
		console.log('options.setSortBy', sortBy)
		d.push('setSortBy', sortBy)
	},
	
	getOptions: function() {
		d.push('setOptions')
	},

	// Initializer
	toProperty: function(initialOptions) {
		console.log('options-before', initialOptions)
		// const gotOptions = 
		// 	.scan(initialOptions, (_, newOptions) => newOptions)
		// 	.log('options')

		return Bacon.update(
			initialOptions,
			[d.stream('setSortBy')], setSortBy,
			[d.stream('setOptions')], (_, newOptions) => newOptions
		)
		.log('options.baconUpdate')

		function setSortBy(options, sortBy) {
			return R.merge(options, {sortBy: sortBy, firstRun: false})
		}

		function doSetOptions(options, newOptions) {
			return options
		}
	}
}

