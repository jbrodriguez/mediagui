const Bacon = require('baconjs'),
	R = require('ramda'),
	Dispatcher = require('./dispatcher'),
	api = require('./api')

const d = new Dispatcher()

module.exports = {
	toProperty: function(initialMovies, optionS) {
		gotMovies = d.stream('getMovies')
			.flatMap(options => Bacon.fromPromise(api.getMovies(options)))

		const itemsS = Bacon.update(
			initialMovies,
			gotMovies, (movies, newMovies) => movies
		)

		return Bacon.combineAsArray([itemsS, optionsS]).map(withDisplayStatus)

		function movies(items, newItems) {
			return newItems
		}
	},

	getOptions: function(options) {
		d.push('getOptions', options)
	}

}

