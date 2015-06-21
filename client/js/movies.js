const Bacon       = require('baconjs'),
      R           = require('ramda'),
      Dispatcher  = require('./dispatcher'),
      api         = require('./api')


const d = new Dispatcher()

module.exports = {

    // "public" methods
    getCover: function() {
    	console.log('movies.getCover')
    	d.push('getCover')
    },

    getMovies: function(options) {
    	console.log('movies.getMovies: ', options)
        d.push('getMovies', options)
    },

    importMovies: function() {
    	d.push('importMovies')
    },

    setMovieScore: function(movie, score) {
    	movie.score = score
    	console.log("movies.setMovieScore: ", JSON.stringify(movie))
    	d.push('setScore', movie)
    },

    toProperty: function(initialMovies, optionsS) {
    	console.log('movies-before')
        const gotMovies = d
        	.stream('getMovies')
        	.log('movies-opt')
            .flatMap(opt => Bacon.fromPromise(api.getMovies(opt)))
            .log('movies-middle')

        const gotCover = d
        	.stream('getCover')
        	.flatMap(_ => Bacon.fromPromise(api.getCover()))
        	.log('cover')

        const movieImported = d
        	.stream('importMovies')
        	.flatMap(_ => Bacon.fromPromise(api.importMovies()))
        	.log('importMovies')

        const movieScoreChanged = d
        	.stream('setScore')
        	.flatMap( (movie) => Bacon.fromPromise( api.setMovieScore(movie) ))

        optionsS.onValue((opt) => {
        	console.log('movies.optionsS.onValue', opt)
        	if (!opt.firstRun) {
	        	this.getMovies(opt)
	        }
        })

        return Bacon.update(
        	initialMovies,
        	[gotMovies], (_, newMovies) => newMovies,
        	[gotCover], (_, newCover) => newCover,
        	[movieImported], (currentMovies, _) => currentMovies,
        	movieScoreChanged, doMovieScoreChanged
        )
        .log('movies')

        function doMovieScoreChanged(movies, changedMovie) {
        	var id = changedMovie.id,
        		score = changedMovie.score

        	const items = R.map(updateItem(id, it => R.merge(it, {score})), movies.items)
        	return R.merge(movies, {items})
        }
    }

}

function updateItem(itemId, fn) {
	return (it) => it.id === itemId ? fn(it) : it
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