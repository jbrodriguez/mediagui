// const Bacon       = require('baconjs'),
//       R           = require('ramda'),
const 	fetch 		= require('jquery').ajax,
		// websocket 	= require('./websocket'),
		Bacon 		= require('baconjs')

const hostr	= "http://" + document.location.host + "/api/v1"
const hostw = "ws://" + document.location.host + "/ws"

module.exports = {
	getSocket: getSocket,
	getConfig: getConfig,
	getCover: getCover,
	getMovies: getMovies,
	importMovies: importMovies,
	addMediaFolder: addMediaFolder,
	setMovieScore: setMovieScore,
	setMovieWatched: setMovieWatched,
	fixMovie: fixMovie,
	getDuplicates: getDuplicates,
	pruneMovies: pruneMovies
}

function getSocket() {
	const skt = new WebSocket(hostw)

	skt.onopen = function() {
	    console.log("Connection opened")
	}

	skt.onclose = function() {
	    console.log("Connection is closed...")
	}

	const stream = Bacon.fromEventTarget(skt, "message").map(function(event) {
		// console.log('event is: ', event)
	    var dataString = event.data
	    // console.log("got:", JSON.parse(dataString))
	    return JSON.parse(dataString)
	})

	const sendMsg = function(topic, msg) {
		const message = {
			topic: topic,
			payload: JSON.stringify(msg)
		}

		skt.send(JSON.stringify(message))
	}	

	return {
		socketS: stream,
		sendFn: sendMsg
	}
}

function getConfig() {
	// console.log('inside api.getConfig')

	return fetch(hostr + '/config')
}

function getCover() {
	// console.log('inside api.getCover')

	return fetch(hostr + '/movies/cover')
}

function getMovies(options) {
	return fetch(hostr + '/movies', {
		data: options
	})
}

function importMovies() {
	return fetch(hostr + '/import', {
		method: 'POST'
	})
}

function addMediaFolder(folder) {
	return fetch(hostr + '/config/folder', {
		method: 'PUT',
		data: JSON.stringify({topic: "", payload: folder})
	})
}

function setMovieScore(movie) {
	// console.log("api.setMovieScore: (" + JSON.stringify(movie) + ")")
	return fetch(hostr + '/movies/' + movie.id + '/score', {
		method: 'PUT',
		data: JSON.stringify(movie)
	})
}

function setMovieWatched(movie) {
	// console.log("api.setMovie: (" + JSON.stringify(movie) + ")")
	return fetch(hostr + '/movies/' + movie.id + '/watched', {
		method: 'PUT',
		data: JSON.stringify(movie)
	})
}

function fixMovie(movie) {
	return fetch(hostr + '/movies/' + movie.id + '/fix', {
		method: 'PUT',
		data: JSON.stringify(movie)
	})	
}

function getDuplicates() {
	return fetch(hostr + '/movies/duplicates')
}

function pruneMovies() {
	return fetch(hostr + '/prune', {
		method: 'POST'
	})
}