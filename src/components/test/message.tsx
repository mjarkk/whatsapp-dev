import { Button } from "@/components/ui/button"
import { getUrl } from "@/services/fetch"
import { State } from "@/services/state"
import { FormEvent, useState } from "react"
import { Input } from "@/components/ui/input"
import { Label } from "../ui/label"
import { toast } from "sonner"
import { CodeBlock } from "./codeBlock"

export interface HelloWorldMessageProps {
	state: State
}

export function HelloWorldMessage({ state }: HelloWorldMessageProps) {
	const [phoneNumber, setPhoneNumber] = useState("")
	const [message, setMessage] = useState("Hello world!")

	const url = getUrl(`/v18.0/${state.phoneNumberID}/messages`)
	const auth = `Bearer ${state.graphToken}`
	const body = JSON.stringify({
		messaging_product: "whatsapp",
		to: phoneNumber,
		type: "text",
		text: {
			body: message,
		},
	})

	const onSubmit = async (e: FormEvent<HTMLFormElement>) => {
		e.preventDefault()

		const resp = await fetch(url, {
			method: "POST",
			headers: {
				Authorization: auth,
				"Content-Type": "application/json",
			},
			body,
		})
		if (resp.status >= 400) {
			toast.error(await resp.text())
		} else {
			toast.success("Message sent check the chats below")
		}
	}

	const code = `curl -i -X POST \\
    ${url} \\
    -H 'Authorization: ${auth}' \\
    -H 'Content-Type: application/json' \\
    -d '${body}'`

	return (
		<form onSubmit={onSubmit} w-full>
			<h4>Message</h4>
			<div my-4>
				<Label htmlFor="phoneNumber">Phone number</Label>
				<Input
					name="phoneNumber"
					id="phoneNumber"
					placeholder="31612345678"
					value={phoneNumber}
					onChange={(e) => setPhoneNumber(e.target.value)}
				/>
			</div>
			<div my-4>
				<Label htmlFor="message">Message</Label>
				<Input
					name="message"
					id="message"
					placeholder="Hello world!"
					value={message}
					onChange={(e) => setMessage(e.target.value)}
				/>
			</div>

			<CodeBlock code={code} />

			<Button mt-3 type="submit">
				Send
			</Button>
		</form>
	)
}
