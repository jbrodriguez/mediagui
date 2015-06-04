const	React 			= require('react'),
		Bacon 			= require('baconjs'),
		MediaGUI 		= require('./MediaGUI.jsx'),
		MoviesCover 	= require('./MoviesCover.jsx'),
		MoviesPage 		= require('./MoviesPage.jsx'),
		settings 		= require('./settings'),
		movies 			= require('./movies'),
		api 			= require('./api'),
		storage			= require('./storage'),
		options 		= require('./options'),
		Router 			= require('react-router'),
		Route 			= Router.Route,
		DefaultRoute 	= Router.DefaultRoute		

const	routes 		= (
			<Route name="app" path="/" handler={MediaGUI}>
				<DefaultRoute handler={MoviesCover} />
				<Route name="movies" path="/movies" handler={MoviesPage} />
			</Route>
		)

// const	settingsP 	= settings.toProperty({mediaFolders:[], version:"0.4.0-7.fbb280b"}),
const	settingsP 	= settings.toProperty({}),
		optionsP 	= options.toProperty(getInitialOptions()),
      	moviesP  	= movies.toProperty([])

const	appState 	= Bacon.combineTemplate({
			settings: settingsP,
			movies: moviesP,
			options: optionsP
		})
		.log('hey dude')

appState.onValue((state) => {
	console.log('inventando: ', state.settings)
	if (state.settings != null) {
		console.log('rendering')
		Router.run(routes, Router.HistoryLocation, function(Handler) {
			React.render(<Handler { ...state} />, document.body, function() {
				console.log('marrano')
			})
		})
	}
	console.log('tonight is what it means to be young')
})

// appState.onValue(function(state) {
// 	console.log('inventando: ', state)
// 	React.render(<mediaGUI {...state} />, document.getElementById('app'))
// })

settings.getConfig()

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
        searchTerm: searchTerm,
        filterByOptions: filterByOptions,
        filterBy: filterBy,
        sortByOptions: sortByOptions,
        sortBy: sortBy,
        sortOrderOptions: sortOrderOptions,
        sortOrder: sortOrder,
        mode: mode
	}

	return base
}



// window.onload = function() {
//   React.render(<MediaGui />, document.getElementById('app'))
// }