const	React 		= require('react'),
		Bacon 		= require('baconjs'),
		MediaGUI 	= require('./mediaGUI.jsx'),
		settings 	= require('./settings'),
		movies 		= require('./movies'),
		api 		= require('./api'),
		options 	= require('./options')

// const	settingsP 	= settings.toProperty({mediaFolders:[], version:"0.4.0-7.fbb280b"}),
const	settingsP 	= settings.toProperty({}),
		optionsP 	= options.toProperty({}),
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
		React.render(<MediaGUI { ...state} />, document.getElementById('app'), function() {
			console.log('marrano')
		})
	}
	console.log('tonight is what it means to be young')
})

// appState.onValue(function(state) {
// 	console.log('inventando: ', state)
// 	React.render(<mediaGUI {...state} />, document.getElementById('app'))
// })

settings.getConfig()




// window.onload = function() {
//   React.render(<MediaGui />, document.getElementById('app'))
// }