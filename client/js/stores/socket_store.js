import Bacon from 'baconjs'
import ffux from 'ffux';

const Socket = ffux.createStore({
	actions: ["send"],
	state: (initalMessages, {send}, {wsapi}) => {
		const receiveS = wsapi.get()
		const sendS = send.doAction(wsapi.send)

		return Bacon.update(
			initalMessages,
			receiveS, (local, remote) => local.concat(remote)
		)
	}
})

export default Socket