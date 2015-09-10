import fetch from 'node-fetch'
import Bacon from 'baconjs'
// import USP from 'url-search-params'

export default class Api {
	constructor() {
		this.hostr = "http://" + document.location.host + "/api/v1"
	}

	getConfig() {
		return fetch(this.hostr + '/config')
			.then(resp => resp.json())
	}

	getCover() {
		// console.log('api.getCover')
		return fetch(this.hostr + '/movies/cover')
			.then(resp => resp.json())
	}

	getMovies(options) {
		// console.log('options: ', options)
		var u = toUrl(options)
		// console.log('getMovies.u', u)
		return fetch(this.hostr + '/movies?' + u)
			.then(resp => resp.json())
	}

	importMovies() {
		return fetch(this.hostr + '/import', {method: 'POST'})
			.then(resp => resp.json())
	}

	addMediaFolder(folder) {
		return fetch(this.hostr + '/config/folder', {
			method: 'PUT',
			body: JSON.stringify({topic: "", payload: folder})
		})
		.then(resp => resp.json())
	}

	setMovieScore(movie) {
		// console.log("api.setMovieScore: (" + JSON.stringify(movie) + ")")
		return fetch(this.hostr + '/movies/' + movie.id + '/score', {
			method: 'PUT',
			body: JSON.stringify(movie)
		})
		.then(resp => resp.json())
	}

	setMovieWatched(movie) {
		// console.log("api.setMovie: (" + JSON.stringify(movie) + ")")
		return fetch(this.hostr + '/movies/' + movie.id + '/watched', {
			method: 'PUT',
			body: JSON.stringify(movie)
		})
		.then(resp => resp.json())
	}

	fixMovie(movie) {
		return fetch(this.hostr + '/movies/' + movie.id + '/fix', {
			method: 'PUT',
			body: JSON.stringify(movie)
		})	
		.then(resp => resp.json())
	}

	getDuplicates() {
		return fetch(this.hostr + '/movies/duplicates')
			.then(resp => resp.json())
	}

	pruneMovies() {
		return fetch(this.hostr + '/prune', { method: 'POST' })
			.then(resp => resp.json())
	}	

}

function toUrl(dict) {
	var query = [], i, key, name, value

	for (key in dict) {
		name = encode(key)
		value = encode(dict[key])
	
		query.push(name + '=' + value)
	}

	return query.join('&')
}

function encode(str) {
	return encodeURIComponent(str).replace(find, replacer);
}

const 
	find = /[!'\(\)~]|%20|%00/g,
	replace = {
		'!': '%21',
		"'": '%27',
		'(': '%28',
		')': '%29',
		'~': '%7E',
		'%20': '+',
		'%00': '\x00'
	},
	replacer = function(match) {
		return replace[match];
	}

// const Bacon       = require('baconjs'),
//       R           = require('ramda'),
// const 	fetch 		= require('jquery').ajax,
// 		// websocket 	= require('./websocket'),
// 		Bacon 		= require('baconjs')

// const hostr	= "http://" + document.location.host + "/api/v1"

// module.exports = {
// 	// getSocket: getSocket,
// 	getConfig: getConfig,
// 	getCover: getCover,
// 	getMovies: getMovies,
// 	importMovies: importMovies,
// 	addMediaFolder: addMediaFolder,
// 	setMovieScore: setMovieScore,
// 	setMovieWatched: setMovieWatched,
// 	fixMovie: fixMovie,
// 	getDuplicates: getDuplicates,
// 	pruneMovies: pruneMovies
// }

// // function getSocket() {
// // 	const skt = new WebSocket(hostw)

// // 	skt.onopen = function() {
// // 	    console.log("Connection opened")
// // 	}

// // 	skt.onclose = function() {
// // 	    console.log("Connection is closed...")
// // 	}

// // 	const stream = Bacon.fromEventTarget(skt, "message").map(function(event) {
// // 		// console.log('event is: ', event)
// // 	    var dataString = event.data
// // 	    // console.log("got:", JSON.parse(dataString))
// // 	    return JSON.parse(dataString)
// // 	})

// // 	const sendMsg = function(topic, msg) {
// // 		const message = {
// // 			topic: topic,
// // 			payload: JSON.stringify(msg)
// // 		}

// // 		skt.send(JSON.stringify(message))
// // 	}	

// // 	return {
// // 		socketS: stream,
// // 		sendFn: sendMsg
// // 	}
// // }

// function getConfig() {
// 	// console.log('inside api.getConfig')

// 	return fetch(hostr + '/config')
// }

// function getCover() {
// 	console.log('inside api.getCover')

// 	return fetch(hostr + '/movies/cover')
// }

// function getMovies(options) {
// 	console.log('opt: ', toUrl(options))
// 	return fetch(hostr + '/movies', {
// 		data: options
// 	})
// }

// function importMovies() {
// 	return fetch(hostr + '/import', {
// 		method: 'POST'
// 	})
// }

// function addMediaFolder(folder) {
// 	return fetch(hostr + '/config/folder', {
// 		method: 'PUT',
// 		data: JSON.stringify({topic: "", payload: folder})
// 	})
// }

// function setMovieScore(movie) {
// 	// console.log("api.setMovieScore: (" + JSON.stringify(movie) + ")")
// 	return fetch(hostr + '/movies/' + movie.id + '/score', {
// 		method: 'PUT',
// 		data: JSON.stringify(movie)
// 	})
// }

// function setMovieWatched(movie) {
// 	// console.log("api.setMovie: (" + JSON.stringify(movie) + ")")
// 	return fetch(hostr + '/movies/' + movie.id + '/watched', {
// 		method: 'PUT',
// 		data: JSON.stringify(movie)
// 	})
// }

// function fixMovie(movie) {
// 	return fetch(hostr + '/movies/' + movie.id + '/fix', {
// 		method: 'PUT',
// 		data: JSON.stringify(movie)
// 	})	
// }

// function getDuplicates() {
// 	return fetch(hostr + '/movies/duplicates')
// }

// function pruneMovies() {
// 	return fetch(hostr + '/prune', {
// 		method: 'POST'
// 	})
// }