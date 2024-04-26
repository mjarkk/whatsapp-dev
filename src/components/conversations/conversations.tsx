import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { fetch, post } from "@/services/fetch"
import { FormEvent, useEffect, useRef, useState } from "react"
import { OpenCloseButton } from "../openCloseButton"
import { useConversationsStore, type Conversation } from "@/services/state"
import { NewChatDialog } from "./new"
import { ShowMessage } from "./singleMessage"

export function Conversations() {
	const [open, setOpen] = useState(true)
	const [newConversationOpen, setNewConversationOpen] = useState(false)
	const { conversations, setConversations, newConversation } =
		useConversationsStore()

	const getData = async () => {
		const conversationsResponse = await fetch("/api/conversations")
		setConversations(await conversationsResponse.json())
	}

	useEffect(() => {
		getData()
	}, [])

	return (
		<>
			<h2 m-6 mb-0 flex flex-wrap gap-4 justify-between items-center>
				<span inline-flex items-center>
					<OpenCloseButton open={open} setOpen={setOpen} /> Conversations
				</span>
				<Button onClick={() => setNewConversationOpen(true)}>
					New conversation!
				</Button>
			</h2>
			{conversations && open ? (
				<div flex flex-wrap gap-4 p-4>
					{conversations.map((conversation) => (
						<Conversation
							conversation={conversation}
							key={conversation.phoneNumberId}
						/>
					))}
				</div>
			) : undefined}
			<NewChatDialog
				open={newConversationOpen}
				newConversation={newConversation}
				close={() => setNewConversationOpen(false)}
			/>
		</>
	)
}

interface ConversationProps {
	conversation: Conversation
}

function Conversation({ conversation }: ConversationProps) {
	const { updateConversation } = useConversationsStore()
	const [msgCount, setMsgCount] = useState(0)
	const messagesEndRef = useRef<HTMLDivElement>(null)

	const onSendMessage = async (e: FormEvent<HTMLFormElement>) => {
		e.preventDefault()

		const target = e.target as HTMLFormElement

		const formData = new FormData(target)
		const message = Object.fromEntries(formData).message
		if (message === "") return

		target.reset()
		const response = await post(`/api/conversations/${conversation.ID}`, {
			message,
		})
		updateConversation(await response.json())
	}

	useEffect(() => {
		if (msgCount === conversation.messages.length) {
			return
		}

		const smooth = msgCount > 0
		setMsgCount(conversation.messages.length)

		if (messagesEndRef.current) {
			messagesEndRef.current.scrollIntoView({
				behavior: smooth ? "smooth" : "instant",
			})
		}
	}, [msgCount, conversation])

	return (
		<div key={conversation.phoneNumber} bg-zinc-900 w-100 rounded>
			<h4
				m-0
				p-3
				border-solid
				border-0
				border-b-2
				border-zinc-700
				text-zinc-200
			>
				{conversation.phoneNumber}
			</h4>
			<div h-130 overflow-y-auto>
				<div flex flex-col justify-end>
					{conversation.messages.map((message) => (
						<ShowMessage key={message.whatsappID} message={message} />
					))}
					<div ref={messagesEndRef} />
				</div>
			</div>
			<form bg-zinc-700 onSubmit={onSendMessage}>
				<Input type="text" name="message" placeholder="message" />
			</form>
		</div>
	)
}
