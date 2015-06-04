// const Bacon       = require('baconjs'),
//       R           = require('ramda'),
const 	fetch 	= require('jquery')

const api = "http://localhost:7623/api/v1"

function getConfig() {
	console.log('inside api.getConfig')
	const sup = fetch.ajax(api + '/config')
		// .then(function(res) {
	 //        return res.json();
	 //    })
	 //    .then(function(json) {
	 //    	console.log("rocky: ", json)
	 //    	return json
	 //    })

	console.log('never surrender', sup)
	return sup
}

function getCover() {
	console.log('inside api.getCover')
	return fetch.ajax(api + '/movies/cover')
}

function getMovies(options) {
	return fetch(api + '/movies')
}

module.exports = {
	getConfig: getConfig,
	getCover: getCover,
	getMovies: getMovies
}