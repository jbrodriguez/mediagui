const	React 			= require('react'),
		Bacon 			= require('baconjs'),
		MediaGUI 		= require('./MediaGUI.jsx'),
		MoviesCover 	= require('./MoviesCover.jsx'),
		MoviesPage 		= require('./MoviesPage.jsx'),
		SettingsPage	= require('./SettingsPage.jsx'),
		Import 			= require('./Import.jsx'),
		PrunePage 		= require('./PrunePage.jsx'),
		settings 		= require('./settings'),
		movies 			= require('./movies'),
		wsmessages		= require('./wsmessages'),
		Dispatcher  	= require('./dispatcher'),
		api 			= require('./api'),
		storage			= require('./storage'),
		options 		= require('./options'),
		Router 			= require('react-router'),
		Route 			= Router.Route,
		DefaultRoute 	= Router.DefaultRoute,
		Redirect 		= Router.Redirect

const d = new Dispatcher()

api
	.getConfig()
	.then(function(config) {
		run(config)
	})

function run(config) {
	console.log(config)

 	const {socketS, sendFn} = api.getSocket()

	const	navigationS	= d.stream('navigation'),
			settingsP 	= settings.toProperty(config),
			optionsP 	= options.toProperty(getInitialOptions()),
		  	moviesP  	= movies.toProperty({}, optionsP),
		  	messageP 	= wsmessages.toProperty([], socketS, sendFn)

	const	appState 	= Bacon.combineTemplate({
				settings: settingsP,
				movies: moviesP,
				options: optionsP,
				messages: messageP,
				navigation: navigationS
			})
			// .log('appState.value = ')

	const	routes 		= (
				<Route name="app" path="/" handler={MediaGUI}>
					<Route name="cover" path="/movies/cover" handler={MoviesCover} />
					<Route name="movies" path="/movies" handler={MoviesPage} />
					<Route name="settings" path="/settings" handler={SettingsPage} />
					<Route name="import" path="/import" handler={Import} />
					<Route name="duplicates" path="/movies/duplicates" handler={MoviesPage} />
					<Route name="prune" path="/prune" handler={PrunePage} />

					<Redirect from="/" to="/movies/cover" />
				</Route>
			)

	var Handler = {}

	const router = Router.create({
		routes: routes,
		location: Router.HistoryLocation
	})

	if (config.mediaFolders.length == 0) {
		// console.log("should have piaid me")
		router.transitionTo("settings")
	}
	// } else {
	// 	router.transitionTo("cover")
	// }

	router.run( function(ProxyHandler, state) {
		Handler = ProxyHandler

		// console.log('router.run.state: ', state)

		if (state.routes.length > 1) {
			var minus1 = state.routes.length - 1
			// console.log('state.routes['+minus1+'].path = ' + state.routes[state.routes.length - 1].path)

			switch (state.routes[minus1].path) {
				case "/movies":
					movies.getMovies(state.query)
					break;
				case "/movies/cover":
					movies.getCover()
					break;
				case "/movies/duplicates":
					movies.getDuplicates()
					break;
				case "/import":
				case "/settings":
				case "/prune":
					d.push('navigation')
					break;
			}
		}
	})

	appState.onValue((state) => {
		// console.log('dentro de onValue: ', state)
		React.render(<Handler { ...state}/>, document.body)
		// React.render(<Handler { ...state} />, document.body, function() {
		// 	console.log('marrano')
		// })
	})

	d.push('navigation')	
}

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