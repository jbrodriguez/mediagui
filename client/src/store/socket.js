class Socket {
	constructor() {
		const host = `ws://${document.location.host}/ws`
		this.skt = new WebSocket(host)
		this.skt.onopen = () => console.log('Connection opened') // eslint-disable-line
		this.skt.onclose = () => console.log('Connection is closed...') // eslint-disable-line
	}

	receive(fn) {
		this.skt.onmessage = fn
	}

	send({ topic, payload }) {
		const packet = {
			topic,
			payload: JSON.stringify(payload),
		}

		this.skt.send(JSON.stringify(packet))
	}
}

export default new Socket()
