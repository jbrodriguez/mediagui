import React from 'react'
import ffux from 'ffux'
// import { Router, Route, Redirect } from 'react-router'
import { Listener } from 'ffux/react'

import Api from './lib/api.js'
import WebsocketApi from './lib/wsapi.js'
import storage from './lib/storage.js'

import Movies from './stores/movies_store'
import Settings from './stores/settings_store'
import Options from './stores/options_store'
import Socket from './stores/socket_store'

import App from './screens/app_page.js'
import CoversPage from './screens/covers_page'
import MoviesPage from './screens/movies_page'
import SettingsPage from './screens/settings_page'
import DuplicatesPage from './screens/duplicates_page'
import ImportPage from './screens/import_page'
import PrunePage from './screens/prune_page'

const	Router 			= require('react-router'),
		Route 			= Router.Route,
		Redirect 		= Router.Redirect

const wsapi = new WebsocketApi()
const api = new Api()

api
	.getConfig()
	.then(function(config) {
		run(config)
	})

function run(config) {
	const 	settings = Settings(config, {api}),
			options = Options(getInitialOptions(), {storage}),
			movies = Movies({}, {api, options}),
			messages = Socket([], {wsapi})



	class Frame extends React.Component {
		render() {
			return (
				<Listener 
						initialState={{}}
						dispatcher={ _ => ffux({
							settings,
							options,
							movies,
							messages,
						})}>
					<App />
				</Listener>
			)
		}
	}

	const routes = (
		<Route name="app" path="/" handler={Frame}>
			<Route name="cover" path="/cover" handler={CoversPage} />
			<Route name="movies" path="/movies" handler={MoviesPage} />
			<Route name="settings" path="/settings" handler={SettingsPage} />
			<Route name="import" path="/import" handler={ImportPage} />
			<Route name="duplicates" path="/duplicates" handler={DuplicatesPage} />
			<Route name="prune" path="/prune" handler={PrunePage} />

			<Redirect from="/" to="/cover" />
		</Route>
	)

	const router = Router.create({
		routes: routes,
		location: Router.HistoryLocation
	})

	if (config.mediaFolders.length == 0) {
		router.transitionTo("settings")
	}

	router.run(function(Handler) {
		React.render(<Handler />, document.body)
	})
}

function getInitialOptions() {
	var searchTerm = ''

    var filterByOptions = [
        {value: 'title', label: 'Title'}, 
        {value: 'genre', label: 'Genre'},
        {value: 'year', label: 'Year'},
        {value: 'country', label: 'Country'},
        {value: 'director', label: 'Director'},
        {value: 'actor', label: 'Actor'}
    ]
    var filterBy = storage.get('filterBy') || 'title'

    var sortByOptions = [
        {value: 'title', label: 'Title'}, 
        {value: 'runtime', label: 'Runtime'}, 
        {value: 'added', label: 'Added'}, 
        {value: 'last_watched', label: 'Watched W'}, 
        {value: 'count_watched', label: 'Watched C'}, 
        {value: 'year', label: 'Year'}, 
        {value: 'imdb_rating', label: 'Rating'}
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
        offset: 0
	}

	return base
}