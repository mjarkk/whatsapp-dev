import { fetch } from "@/services/fetch"
import { Fragment, useEffect, useState } from "react"
import { Conversations } from "@/components/conversations/conversations"
import { Templates } from "@/components/templates/templates"
import { Test } from "@/components/test/test"
import { State, useConversationsStore } from "@/services/state"
import { EventsWebsocket } from "@/services/websocket"

export function App() {
	const [state, setState] = useState<State>({
		graphToken: "",
		appSecret: "",
		phoneNumber: "",
		phoneNumberID: "",
		webhookURL: "",
	})

	const getData = async () => {
		const stateResponse = await fetch("/api/info")
		setState(await stateResponse.json())
	}

	useEffect(() => {
		getData()
	}, [])

	return (
		<div>
			<div p-5 bg-zinc-900>
				<h1 m-0>Whatsapp Dev</h1>
				<p m-0>
					{state.phoneNumber}{" "}
					<span italic text-zinc-400>
						(id: {state.phoneNumberID})
					</span>
				</p>
			</div>

			<Test state={state} />

			<Templates />

			<Conversations />

			<WebsocketHandler />
		</div>
	)
}

function WebsocketHandler() {
	const { addMessage } = useConversationsStore()

	useEffect(() => {
		const ws = new EventsWebsocket((data) => {
			console.log("websocket message:", data)
			if (data.type === "message") {
				addMessage(data.message)
			}
		})
		ws.start()
		return () => ws.close()
	}, [])

	return <Fragment></Fragment>
}
