import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { post } from "@/services/fetch"
import { useState } from "react"
import type { Conversation } from "@/services/state"

export interface NewChatDialogProps {
	open: boolean
	newConversation: (conversation: Conversation) => void
	close: () => void
}

export function NewChatDialog({
	newConversation,
	open,
	close,
}: NewChatDialogProps) {
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
