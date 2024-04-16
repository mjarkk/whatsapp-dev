import { getUrl } from "./fetch"

export class EventsWebsocket {
	private closed = false
	private socket: WebSocket | null = null
	private onEvent: (data: any) => void

	constructor(onEvent: (data: any) => void) {
		this.onEvent = onEvent
	}

	async start() {
		if (this.socket) {
			throw "cannot start a Websocket twice"
		}

		while (true) {
			await new Promise((res) => {
				this.socket = new WebSocket(getUrl("/api/events", true))
				this.socket.onmessage = (ev) => {
					if (this.closed) {
						return
					}
					this.onEvent(JSON.parse(ev.data))
				}
				this.socket.onclose = () => setTimeout(res, 5_000)
			})
			if (this.closed) {
				break
			}
		}
	}

	close() {
		this.closed = true
		this.socket?.close()
	}
}
