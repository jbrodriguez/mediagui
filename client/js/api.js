// const Bacon       = require('baconjs'),
//       R           = require('ramda'),
const 	fetch 	= require('jquery')

function getConfig() {
	console.log('inside api')
	const sup = fetch.ajax('http://localhost:7623/api/v1/config')
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

function getMovies(options) {
	return fetch('/movies')
}

module.exports = {
	getConfig: getConfig,
	getMovies: getMovies
}