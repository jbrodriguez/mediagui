const	React   = require('react'),
		// Bacon = require('baconjs'),
		MediaGui = require('./mediaGUI.jsx')
// 		movies   = require('./movies'),
// 		options  = require('./options')

// const optionsP = options.toProperty({}),
//       moviesP  = movies.toProperty([], optionsP)

// const appState = Bacon.combineTemplate({
//   movies: moviesP,
//   options: optionsP
// })

// appState.onValue((state) => {
//   React.render(<mediaGUI {...state} />, document.getElementById('app'))
// })

window.onload = function() {
  React.render(<MediaGui />, document.getElementById('app'))
}