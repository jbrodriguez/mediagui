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

	setOffset: function(offset) {
		console.log('options.setOffset', offset)
		d.push('setOffset', offset)
	},

	setQueryTerm: function(term) {
		console.log('options.setQueryTerm', term)
		d.push('setQueryTerm', term)
	},

	// Initializer
	toProperty: function(initialOptions) {
		console.log('options-before', initialOptions)

		const gotQueryTerm = d
			.stream('setQueryTerm')
			.debounce(750)

		// const gotOptions = 
		// 	.scan(initialOptions, (_, newOptions) => newOptions)
		// 	.log('options')

		return Bacon.update(
			initialOptions,
			[d.stream('setSortBy')], setSortBy,
			[d.stream('setSortOrder')], setSortOrder,
			[d.stream('setOffset')], setOffset,
			gotQueryTerm, setQueryTerm
		)
		.log('options.baconUpdate')

		function setSortBy(options, sortBy) {
			return R.merge(options, {sortBy: sortBy, firstRun: false})
		}

		function setSortOrder(options, sortOrder) {
			return R.merge(options, {sortOrder: sortOrder, firstRun: false})
		}

		function setOffset(options, offset) {
			return R.merge(options, {offset: offset, firstRun: false})
		}

		function setQueryTerm(options, term) {
			return R.merge(options, {query: term, firstRun: false})
		}
	}
}

