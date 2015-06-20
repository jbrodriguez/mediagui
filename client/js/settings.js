const Bacon 		= require('baconjs'),
	  R 			= require('ramda'),
	  Dispatcher 	= require('./dispatcher'),
	  api 			= require('./api')

const d = new Dispatcher()

module.exports = {
	// Public API
	getConfig: function() {
		console.log('getting config')
		d.push('getConfig')
	},

	addMediaFolder: function(folder) {
		d.push('addFolder', folder)
	},

	toProperty: function(initialConfig) {
		console.log('settings-before')
		const gotConfig = d
			.stream('getConfig')
			.flatMap( name => Bacon.fromPromise( api.getConfig() ) )
			.log('settings-middle')

		const addedFolder = d
			.stream('addFolder')
			.flatMap( folder => Bacon.fromPromise( api.addMediaFolder(folder) ) )
			.log('settings-addFolder')


		return Bacon.update(
			initialConfig,
			gotConfig, (_, newConfig) => newConfig,
			addedFolder, doAddFolder
		)
		.log('settings-final:')

		function doAddFolder(config, newFolder) {
			return R.merge(config, {mediaFolders: config.mediaFolders.concat([newFolder])})
		}

	}
}