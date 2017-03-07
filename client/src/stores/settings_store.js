import Bacon from 'baconjs'
import ffux from 'ffux'

const Settings = ffux.createStore({
	actions: ["getConfig", "addMediaFolder"],

	state: (initialSettings, {getConfig, addMediaFolder}, {api}) => {
		const getConfigS = getConfig.flatMap(_ => Bacon.fromPromise(api.getConfig()))
		const addMediaFolderS = addMediaFolder.flatMap(folder => Bacon.fromPromise(api.addMediaFolder(folder)))

		return Bacon.update(
			initialSettings,
			getConfigS, (_, remote) => remote,
			addMediaFolderS, (_, remote) => remote
		)
	}
})

export default Settings