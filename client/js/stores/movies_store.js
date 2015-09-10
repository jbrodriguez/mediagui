import Bacon from 'baconjs'
import R from 'ramda'
import ffux from 'ffux'

const Movies = ffux.createStore({
	actions: ["getCover", "getMovies", "importMovies", "setMovieScore", "setMovieWatched", "fixMovie", "getDuplicates", "pruneMovies"],

	state: (initialMovies, {getCover, getMovies, importMovies, setMovieScore, setMovieWatched, fixMovie, getDuplicates, pruneMovies}, {api, options}) => {
		var _proxy = function(options) {
			const proxy = {
				query: options.query,
				filterBy: options.filterBy,
				sortBy: options.sortBy,
				sortOrder: options.sortOrder,
				limit: options.limit,
				offset: options.offset
			}

			return api.getMovies(proxy)
		}

		// console.log('options: ', options)
		const optionsS = options.toEventStream().skip(1).skipDuplicates()

		// console.log('getMovies: ', getMovies)
		// console.log('optionsS: ', optionsS)
		const moviesS = getMovies.merge(optionsS)

		const getCoverS = getCover.flatMap(Bacon.fromPromise(api.getCover()))
		const getMoviesS = moviesS.flatMap(opt => Bacon.fromPromise(_proxy(opt)))
		const importMoviesS = importMovies.flatMap( _ => Bacon.fromPromise(api.importMovies()))
		const setMovieScoreS = setMovieScore.flatMap(movie => Bacon.fromPromise(api.setMovieScore(movie)))
		const setMovieWatchedS = setMovieWatched.flatMap(movie => Bacon.fromPromise(api.setMovieWatched(movie)))
		const fixMovieS = fixMovie.flatMap(movie => Bacon.fromPromise(api.fixMovie(movie)))
		const getDuplicatesS = getDuplicates.flatMap(_ => Bacon.fromPromise(api.getDuplicates()))
		const pruneMoviesS = pruneMovies.flatMap(_ => Bacon.fromPromise(api.pruneMovies()))
		// const optionsS = options.skip(1).flatMap(opt => Bacon.fromPromise(doProxy(opt)))

		return Bacon.update(
			initialMovies,
			getCoverS, (_, remote) => remote,
			getMoviesS, (_, remote) => remote,
			setMovieScoreS, _setMovieScore,
			setMovieWatchedS, _setMovieWatched,
			fixMovieS, _fixMovie,
			getDuplicatesS, (_, remote) => remote
			// optionsS, (_, remote) => remote
		)

		function _setMovieScore(movies, changedMovie) {
        	var id = changedMovie.id,
        		score = changedMovie.score

        	const items = R.map(updateItem(id, it => R.merge(it, {score})), movies.items)
        	return R.merge(movies, {items})			
		}

		function _setMovieWatched(movies, changedMovie) {
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

		function _fixMovie(movies, changedMovie) {
        	var id = changedMovie.id

        	const items = R.map(updateItem(id, it => R.merge(it, changedMovie)), movies.items)
        	return R.merge(movies, {items})
		}
	}
})

function updateItem(itemId, fn) {
	return (it) => it.id === itemId ? fn(it) : it
}

export default Movies