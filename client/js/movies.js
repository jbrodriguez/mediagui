const Bacon       	= require('baconjs'),
      R           	= require('ramda'),
      Dispatcher  	= require('./dispatcher'),
	  moment 		= require('moment'),
      api         	= require('./api')


const d = new Dispatcher()

module.exports = {

    // "public" methods
    getCover: function() {
    	// console.log('movies.getCover')
    	d.push('getCover')
    },

    getMovies: function(options) {
    	// console.log('movies.getMovies: ', options)
        d.push('getMovies', options)
    },

    importMovies: function() {
    	d.push('importMovies')
    },

    setMovieScore: function(movie, score) {
    	movie.score = score
    	// console.log("movies.setMovieScore: ", JSON.stringify(movie))
    	d.push('setScore', movie)
    },

    setMovieWatched: function(movie, watched) {
    	movie.last_watched = moment.utc(watched).format()
    	// console.log("movies.setMovieWatched: ", JSON.stringify(movie))
    	d.push('setWatched', movie)
    },

    fixMovie: function(movie, tmdb_id) {
    	movie.tmdb_id = tmdb_id
    	d.push('fixMovie', movie)
    },

    getDuplicates: function() {
    	d.push('getDuplicates')
    },

    pruneMovies: function() {
    	d.push('pruneMovies')
    },

    toProperty: function(initialMovies, optionsS) {
    	// console.log('movies-before')
        const gotMovies = d
        	.stream('getMovies')
        	// .log('movies-opt')
            .flatMap(opt => Bacon.fromPromise(api.getMovies(opt)))
            // .log('movies-middle')

        const gotCover = d
        	.stream('getCover')
        	// .log()
        	.flatMap( _ => Bacon.fromPromise( api.getCover() ))
        	.log('cover')

        const movieImported = d
        	.stream('importMovies')
        	.flatMap( (_) => Bacon.fromPromise(api.importMovies()))
        	// .log('importMovies')

        const movieScoreChanged = d
        	.stream('setScore')
        	.flatMap( (movie) => Bacon.fromPromise( api.setMovieScore(movie) ))

        const movieWatchedChanged = d
        	.stream('setWatched')
        	.flatMap( (movie) => Bacon.fromPromise( api.setMovieWatched(movie) ))
        	// .log('movie-watched')

        const movieFixed = d
        	.stream('fixMovie')
        	.flatMap( (movie) => Bacon.fromPromise( api.fixMovie(movie) ))

        const gotDuplicates = d
        	.stream('getDuplicates')
        	.flatMap( _ => Bacon.fromPromise( api.getDuplicates() ))

        const moviesPruned = d
        	.stream('pruneMovies')
        	.flatMap( (_) => Bacon.fromPromise(api.pruneMovies()))
        	// .log('importMovies')        	


        optionsS.onValue((opt) => {
        	// console.log('movies.options', opt)
        	if (!opt.firstRun) {
	        	this.getMovies(opt)
	        }
        })

        return Bacon.update(
        	initialMovies,
        	gotMovies, (_, newMovies) => newMovies,
        	gotCover, (_, coverMovies) => coverMovies,
        	movieImported, (currentMovies, _) => currentMovies,
        	movieScoreChanged, doMovieScoreChanged,
        	movieWatchedChanged, doMovieWatchedChanged,
        	movieFixed, doMovieFixed,
        	gotDuplicates, (_, newMovies) => newMovies,
        	moviesPruned,  (currentMovies, _) => currentMovies
        )


        function doMovieScoreChanged(movies, changedMovie) {
        	var id = changedMovie.id,
        		score = changedMovie.score

        	const items = R.map(updateItem(id, it => R.merge(it, {score})), movies.items)
        	return R.merge(movies, {items})
        }

        function doMovieWatchedChanged(movies, changedMovie) {
        	var id = changedMovie.id,
        		changed = {
	        		last_watched: changedMovie.last_watched,
	        		all_watched: changedMovie.all_watched,
	        		count_watched: changedMovie.count_watched,
	        		modified: changedMovie.modified
	        	}

        	const items = R.map(updateItem(id, it => R.merge(it, changed)), movies.items)
        	return R.merge(movies, {items})
        }        

        function doMovieFixed(movies, changedMovie) {
        	var id = changedMovie.id

        	const items = R.map(updateItem(id, it => R.merge(it, changedMovie)), movies.items)
        	return R.merge(movies, {items})
        }
    }

}

function updateItem(itemId, fn) {
	return (it) => it.id === itemId ? fn(it) : it
}