class Socket {
	private skt: WebSocket
	private ep: string =
		document && document.location
			? `${document.location.protocol === 'http:' ? 'ws:' : 'wss:'}//${document.location.hostname}:7623/ws`
			: `http://blackbeard.apertoire.org:7623/api/v1`

	constructor() {
		this.skt = new WebSocket(this.ep)
		this.skt.onopen = () => console.log('Connection opened') // tslint:disable-line
		this.skt.onclose = () => console.log('Connection is closed...') // tslint:disable-line
	}

	public receive(fn: any) {
		this.skt.onmessage = fn
	}

	public send(msg: any) {
		const packet = {
			topic: msg.topic,
			payload: JSON.stringify(msg.payload),
		}

		this.skt.send(JSON.stringify(packet))
	}
}

export default new Socket()
