import fetch from 'node-fetch'
import Bacon from 'baconjs'

export default class Api {
	constructor() {
		this.hostr = "http://" + document.location.host + "/api/v1"
	}

	getConfig() {
		return fetch(this.hostr + '/config')
			.then(resp => resp.json())
	}

	getCover() {
		return fetch(this.hostr + '/movies/cover')
			.then(resp => resp.json())
	}

	getMovies(options) {
		var u = toUrl(options)
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
		return fetch(this.hostr + '/movies/' + movie.id + '/score', {
			method: 'PUT',
			body: JSON.stringify(movie)
		})
		.then(resp => resp.json())
	}

	setMovieWatched(movie) {
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