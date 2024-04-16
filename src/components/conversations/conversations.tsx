import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { fetch, post } from "@/services/fetch"
import { FormEvent, useEffect, useState } from "react"
import { OpenCloseButton } from "../openCloseButton"
import {
	useConversationsStore,
	type Conversation,
	type Message,
} from "@/services/state"

export function Conversations() {
	const [open, setOpen] = useState(true)
	const [newConversationOpen, setNewConversationOpen] = useState(false)
	const {
		conversations,
		setConversations,
		newConversation,
		updateConversation,
	} = useConversationsStore()

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
					{conversations.map((conversation, idx) => (
						<Conversation
							conversation={conversation}
							updateConversation={(conv) => updateConversation(idx, conv)}
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

function formatDate(date: Date) {
	const dateFormatter = new Intl.DateTimeFormat("en-US", {
		timeStyle: "short",
	})
	return dateFormatter.format(date)
}

interface ConversationProps {
	conversation: Conversation
	updateConversation: (conversation: Conversation) => void
}

function Conversation({ conversation, updateConversation }: ConversationProps) {
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
		const updatedConverstaion = await response.json()
		updateConversation(updatedConverstaion)
	}

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
			<div h-130 overflow-hidden flex flex-col justify-end>
				{conversation.messages.map((message) => (
					<ShowMessage key={message.whatsappID} message={message} />
				))}
			</div>
			<form bg-zinc-700 onSubmit={onSendMessage}>
				<Input type="text" name="message" placeholder="message" />
			</form>
		</div>
	)
}

function ShowMessage({ message }: { message: Message }) {
	return (
		<div p-2 flex justify={message.direction === "out" ? "start" : "end"}>
			<div
				style={{ maxWidth: "70%" }}
				inline-block
				py-1
				px-2
				bg-zinc-800
				rounded
			>
				{message.headerMessage ? (
					<div font-bold>{message.headerMessage}</div>
				) : undefined}
				<div>
					<span text-xs>{formatDate(new Date(message.timestamp))} - </span>
					{message.message}
				</div>
				{message.footerMessage ? (
					<div font-bold text-sm text-zinc-400 mt-1>
						{message.footerMessage}
					</div>
				) : undefined}
			</div>
		</div>
	)
}

interface NewChatDialogProps {
	open: boolean
	newConversation: (conversation: Conversation) => void
	close: () => void
}

function NewChatDialog({ newConversation, open, close }: NewChatDialogProps) {
	const [state, setState] = useState({
		message: "Hello world!",
		phoneNumber: "",
	})

	const createConversation = async () => {
		const response = await post("/api/conversations", {
			phoneNumber: state.phoneNumber,
			message: state.message,
		})
		const conversation = await response.json()
		newConversation(conversation)
		setState({
			message: "Hello world!",
			phoneNumber: "",
		})
	}

	return (
		<AlertDialog open={open} onOpenChange={() => close()}>
			<AlertDialogContent>
				<AlertDialogHeader>
					<AlertDialogTitle>Create a new conversation</AlertDialogTitle>
					<AlertDialogDescription></AlertDialogDescription>
				</AlertDialogHeader>
				<Label htmlFor="phoneNumber">Source phone number</Label>
				<Input
					value={state.phoneNumber}
					onChange={(e) =>
						setState((s) => ({ ...s, phoneNumber: e.target.value }))
					}
					type="text"
					id="phoneNumber"
					placeholder="+31600000000"
				/>
				<Label htmlFor="message">Message</Label>
				<Input
					value={state.message}
					onChange={(e) => setState((s) => ({ ...s, message: e.target.value }))}
					type="text"
					id="message"
					placeholder="Hello world!"
				/>
				<AlertDialogFooter>
					<AlertDialogCancel>Cancel</AlertDialogCancel>
					<AlertDialogAction onClick={createConversation}>
						Start
					</AlertDialogAction>
				</AlertDialogFooter>
			</AlertDialogContent>
		</AlertDialog>
	)
}
