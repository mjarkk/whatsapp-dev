import {
	MessageButton,
	useConversationsStore,
	type Message,
} from "@/services/state"
import { Button } from "../ui/button"
import { post } from "@/services/fetch"

function formatDate(date: Date) {
	const dateFormatter = new Intl.DateTimeFormat("en-US", {
		timeStyle: "short",
	})
	return dateFormatter.format(date)
}

export function ShowMessage({ message }: { message: Message }) {
	const { updateConversation } = useConversationsStore()

	const buttonReply = async (button: MessageButton) => {
		const response = await post(
			`/api/conversations/${message.conversationId}/btnQuickReply/${button.ID}`,
			{},
		)
		updateConversation(await response.json())
	}

	return (
		<div
			p-2
			flex
			flex-col
			items={message.direction === "out" ? "start" : "end"}
		>
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
			{message.buttons?.length ? (
				<div
					flex
					gap-2
					flex-wrap
					mt-2
					style={{ maxWidth: "70%" }}
					justify="end"
				>
					{message.buttons.map((btn) => (
						<Button
							onClick={() => buttonReply(btn)}
							key={btn.ID}
							variant="secondary"
						>
							{btn.text}
						</Button>
					))}
				</div>
			) : undefined}
		</div>
	)
}
