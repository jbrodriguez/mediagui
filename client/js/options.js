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
	
	setSortOrder: function(sortOrder) {
		console.log('options.setSortOrder', sortOrder)
		d.push('setSortOrder', sortOrder)
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
			[d.stream('setSortOrder')], setSortOrder
		)
		.log('options.baconUpdate')

		function setSortBy(options, sortBy) {
			return R.merge(options, {sortBy: sortBy, firstRun: false})
		}

		function setSortOrder(options, sortOrder) {
			return R.merge(options, {sortOrder: sortOrder, firstRun: false})
		}
	}
}

