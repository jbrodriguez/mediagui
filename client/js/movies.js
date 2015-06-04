const Bacon       = require('baconjs'),
      R           = require('ramda'),
      Dispatcher  = require('./dispatcher'),
      api         = require('./api')


const d = new Dispatcher()

module.exports = {
    toProperty: function(initialMovies) {
    	console.log('movies-before')
        const gotMovies = d
        	.stream('getMovies')
            .flatMap(options => Bacon.fromPromise(api.getMovies(options)))
            .log('movies-middle')

        const gotCover = d
        	.stream('getCover')
        	.flatMap(_ => Bacon.fromPromise(api.getCover()))
        	.log('cover')

        return Bacon.update(
        	initialMovies,
        	[gotMovies], (_, newMovies) => newMovies,
        	[gotCover], (_, newCover) => newCover 
        )
        .log('movies')
    },

    // "public" methods
    getCover: function() {
    	d.push('getCover')
    },

    getMovies: function(options) {
        d.push('getMovies', options)
    }
}

// module.exports = {
//     toProperty: function(initialMovies, optionS) {
//         const gotMovies = d.stream('getMovies')
//                   .flatMap(options => Bacon.fromPromise(api.getMovies(options)))

//         const itemsS = Bacon.update(
//             initialMovies,
//             gotMovies, (movies, newMovies) => movies
//         )

//         return Bacon.combineAsArray([itemsS, filterS]).map(withDisplayStatus)

//         function movies(items, newItems) {
//             return newItems
//         }
//     },

//     // "public" methods
//     getMovies: function(options) {
//         d.push('getMovies', optionS)
//     }
// }