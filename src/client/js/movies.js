const Bacon       = require('baconjs'),
      R           = require('ramda'),
      Dispatcher  = require('./dispatcher')
      api         = require('./api.js')


const d = new Dispatcher()

module.exports = {
    toProperty: function(initialMovies, optionS) {
        gotMovies = d.stream('getMovies')
                  .flatMap(options => Bacon.fromPromise(api.getMovies(options)))

        const itemsS = Bacon.update(
            initialMovies,
            gotMovies, (movies, newMovies) => movies
        )

        return Bacon.combineAsArray([itemsS, filterS]).map(withDisplayStatus)

        function movies(items, newItems) {
            return newItems
        }
    },

    // "public" methods
    getMovies: function(options) {
        d.push('getMovies', optionS)
    }
}