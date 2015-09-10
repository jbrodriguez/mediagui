import Bacon from 'baconjs'

class WebSocketApi {
	constructor() {
		const hostw = "ws://" + document.location.host + "/ws"

		this.skt = new WebSocket(hostw)

		this.skt.onopen = function() {
		    console.log("Connection opened")
		}

		this.skt.onclose = function() {
		    console.log("Connection is closed...")
		}		
	}

	get() {
		const stream = Bacon.fromEventTarget(this.skt, "message").map(function(event) {
			// console.log('event is: ', event)
		    var dataString = event.data
		    // console.log("got:", JSON.parse(dataString))
		    return JSON.parse(dataString)
		})

		return stream		
	}

	send([topic, msg]) {
		const message = {
			topic: topic,
			payload: JSON.stringify(msg)
		}

		this.skt.send(JSON.stringify(message))			
	} 
}

export default WebSocketApi