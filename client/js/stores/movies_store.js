import Bacon from 'baconjs'
import R from 'ramda'
import moment from 'moment'
import ffux from 'ffux'

const Movies = ffux.createStore({
	actions: ["getCover", "getMovies", "importMovies", "setMovieScore", "setMovieWatched", "fixMovie", "getDuplicates", "pruneMovies"],

	state: (initialMovies, {getCover, getMovies, importMovies, setMovieScore, setMovieWatched, fixMovie, getDuplicates, pruneMovies}, {api, options}) => {
		const _options = function(options) {
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

		function _score(movie, score) {
			movie.score = score
			return api.setMovieScore(movie)
		}

		function _watched(movie, watched) {
			movie.last_watched = watched.format()
			return api.setMovieWatched(movie)
		}

		function _fix(movie, tmdb_id) {
			movie.tmdb_id = tmdb_id
			return api.fixMovie(movie)
		}

		const optionsS = options.toEventStream().skip(1).skipDuplicates()
		const moviesS = getMovies.merge(optionsS)

		const getCoverS = getCover.flatMap(Bacon.fromPromise(api.getCover()))
		const getMoviesS = moviesS.flatMap(opt => Bacon.fromPromise(_options(opt)))
		const importMoviesS = importMovies.flatMap(_ => Bacon.fromPromise( api.importMovies() ))
		const setMovieScoreS = setMovieScore.flatMap(([movie, score]) => Bacon.fromPromise(_score(movie, score)))
		const setMovieWatchedS = setMovieWatched.flatMap(([movie, watched]) => Bacon.fromPromise(_watched(movie, watched)))
		const fixMovieS = fixMovie.flatMap(([movie, tmdb_id]) => Bacon.fromPromise(_fix(movie, tmdb_id)))
		const getDuplicatesS = getDuplicates.flatMap(_ => Bacon.fromPromise(api.getDuplicates()))
		const pruneMoviesS = pruneMovies.flatMap(_ => Bacon.fromPromise(api.pruneMovies()))

		return Bacon.update(
			initialMovies,
			getCoverS, (_, remote) => remote,
			getMoviesS, (_, remote) => remote,
			setMovieScoreS, _setMovieScore,
			setMovieWatchedS, _setMovieWatched,
			fixMovieS, _fixMovie,
			getDuplicatesS, (_, remote) => remote,
			importMoviesS, (local, _) => local,
			pruneMoviesS, (local, _) => local
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