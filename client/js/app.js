const	React 			= require('react'),
		Bacon 			= require('baconjs'),
		MediaGUI 		= require('./MediaGUI.jsx'),
		MoviesCover 	= require('./MoviesCover.jsx'),
		MoviesPage 		= require('./MoviesPage.jsx'),
		Import 			= require('./Import.jsx'),
		settings 		= require('./settings'),
		movies 			= require('./movies'),
		Dispatcher  	= require('./dispatcher'),		
		api 			= require('./api'),
		storage			= require('./storage'),
		options 		= require('./options'),
		Router 			= require('react-router'),
		Route 			= Router.Route,
		DefaultRoute 	= Router.DefaultRoute,
		Redirect 		= Router.Redirect

const d = new Dispatcher()

var config = {},
	movieList = []

api
.getConfig()
.then(function(result) {
	config = result
	console.log('obtained getConfig result: ' + config)
	return api.getCover()
})
.then(function(result) {
	movieList = result
	run()
})

function run() {

	// const	settingsP 	= settings.toProperty({mediaFolders:[], version:"0.4.0-7.fbb280b"}),
	const	navigationS	= d.stream('navigation'),
			settingsP 	= settings.toProperty(config),
			optionsP 	= options.toProperty(getInitialOptions()),
		  	moviesP  	= movies.toProperty(movieList, optionsP)

	const	appState 	= Bacon.combineTemplate({
				settings: settingsP,
				movies: moviesP,
				options: optionsP,
				navigation: navigationS
			})
			.log('appState.value = ')

	const	routes 		= (
				<Route name="app" path="/" handler={MediaGUI}>
					<Route name="cover" path="/movies/cover" handler={MoviesCover} />
					<Route name="movies" path="/movies" handler={MoviesPage} />
					<Route name="import" path="/import" handler={Import} />

					<Redirect from="/" to="/movies/cover" />
				</Route>
			)

	var Handler = {}
	
	Router.run(routes, Router.HistoryLocation, function(ProxyHandler, state) {
		Handler = ProxyHandler

		console.log('handler: ', Handler)
		console.log('routes: ' + state.routes)
		console.log('len(routes)=' + state.routes.length)
		console.log('query: ', JSON.stringify(state.query, 4, null))
		if (state.routes.length > 1) {
			var minus1 = state.routes.length - 1
			console.log('state.routes['+minus1+'].path = ' + state.routes[state.routes.length - 1].path)

			switch (state.routes[minus1].path) {
				case "/movies":
					movies.getMovies(state.query)
					break;
				case "/movies/cover":
					movies.getCover()
					break;
				case "/import":
					d.push('navigation')
					break;

			}


		}

		// React.render(<Handler { ...state} />, document.body, function() {
		// 	console.log('marrano')
		// })
	})

	appState.onValue((state) => {
		console.log('dentro de onValue: ', state)
		React.render(<Handler { ...state} />, document.body, function() {
			console.log('marrano')
		})
	})

	d.push('navigation')
}










// appState.onValue(function(state) {
// 	console.log('inventando: ', state)
// 	React.render(<mediaGUI {...state} />, document.getElementById('app'))
// })

// movies.getCover()
// settings.getConfig()

function getInitialOptions() {
	var searchTerm = ''

    var filterByOptions = [
        {id: 1, value: 'title', label: 'Title'}, 
        {id: 2, value: 'genre', label: 'Genre'},
        {id: 3, value: 'country', label: 'Country'},
        {id: 4, value: 'director', label: 'Director'},
        {id: 5, value: 'actor', label: 'Actor'}
    ]
    var filterBy = storage.get('filterBy') || 'title'

    var sortByOptions = [
        {id: 1, value: 'title', label: 'Title'}, 
        {id: 2, value: 'runtime', label: 'Runtime'}, 
        {id: 3, value: 'added', label: 'Added'}, 
        {id: 4, value: 'last_watched', label: 'Watched'}, 
        {id: 5, value: 'year', label: 'Year'}, 
        {id: 6, value: 'imdb_rating', label: 'Rating'}
    ]
    var sortBy = storage.get('sortBy') || 'added'

    var sortOrderOptions = ['asc', 'desc']
    var sortOrder = storage.get('sortOrder') || 'desc'

    var mode = 'regular'

	const base = {
        query: searchTerm,
        filterByOptions: filterByOptions,
        filterBy: filterBy,
        sortByOptions: sortByOptions,
        sortBy: sortBy,
        sortOrderOptions: sortOrderOptions,
        sortOrder: sortOrder,
        mode: mode,
        limit: 50,
        offset: 0,
        firstRun: true
	}

	return base
}



// window.onload = function() {
//   React.render(<MediaGui />, document.getElementById('app'))
// }