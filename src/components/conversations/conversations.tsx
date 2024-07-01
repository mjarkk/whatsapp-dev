import { Button } from "@/components/ui/button"
import { fetch } from "@/services/fetch"
import { useEffect, useState } from "react"
import { OpenCloseButton } from "../openCloseButton"
import { useConversationsStore } from "@/services/state"
import { NewChatDialog } from "./new"
import { Conversation } from "./conversation"

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
