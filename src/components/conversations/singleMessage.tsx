import {
	MessageButton,
	useConversationsStore,
	type Message,
} from "@/services/state"
import { Button } from "../ui/button"
import { post } from "@/services/fetch"
import { text } from "stream/consumers"

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
					<div font-bold>
						<Formatted text={message.headerMessage} />
					</div>
				) : undefined}
				<div>
					<span text-xs>{formatDate(new Date(message.timestamp))} - </span>
					{message.message
						.trim()
						.split("\n")
						.map((line) => line.trim())
						.map((line, index) => (
							<span key={index}>
								{index > 0 ? <br /> : undefined}
								<Formatted text={line} />
							</span>
						))}
				</div>
				{message.footerMessage ? (
					<div font-bold text-sm text-zinc-400 mt-1>
						<Formatted text={message.footerMessage} />
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

function isSpace(c: string): boolean {
	return c === " " || c === "\n"
}

function Formatted({ text }: { text: string }) {
	const parts = [{ bold: false, italic: false, text: "" }]
	for (let idx = 0; idx < text.length; idx++) {
		const c = text[idx]
		const part = parts[parts.length - 1]
		if (c == "_") {
			if (part.italic) {
				// Check if previous part was not a space
				const prevC = text[idx - 1]
				if (!isSpace(prevC) && prevC !== "_") {
					parts.push({ ...part, italic: false, text: "" })
					continue
				}
			} else {
				// Check if the next character is not a space
				const nextC = text[idx + 1]
				if (!isSpace(nextC) && nextC !== "_") {
					parts.push({ ...part, italic: true, text: "" })
					continue
				}
			}
		} else if (c == "*") {
			if (part.bold) {
				// Check if previous part was not a space
				const prevC = text[idx - 1]
				if (!isSpace(prevC) && prevC !== "*") {
					parts.push({ ...part, bold: false, text: "" })
					continue
				}
			} else {
				// Check if the next character is not a space
				const nextC = text[idx + 1]
				if (!isSpace(nextC) && nextC !== "*") {
					parts.push({ ...part, bold: true, text: "" })
					continue
				}
			}
		}
		part.text += c
	}

	const lastPart = parts[parts.length - 1]
	if (lastPart.italic) {
		lastPart.italic = false
		lastPart.text = "_" + lastPart.text
	}
	if (lastPart.bold) {
		lastPart.bold = false
		lastPart.text = "*" + lastPart.text
	}

	if (parts.length > 1) {
		console.log(parts)
	}

	return parts.map((el, idx) => (
		<span
			key={idx}
			style={{
				fontWeight: el.bold ? 600 : 400,
				fontStyle: el.italic ? "italic" : "normal",
			}}
		>
			{el.text}
		</span>
	))
}
