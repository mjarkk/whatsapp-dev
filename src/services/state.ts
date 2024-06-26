import { DBModel } from "@/lib/types"
import { create } from "zustand"

export interface State {
	graphToken: string
	appSecret: string
	phoneNumber: string
	phoneNumberID: string
	webhookURL: string
}

export interface Conversation extends DBModel {
	phoneNumberId: string
	phoneNumber: string
	messages: Array<Message>
}

export interface Message extends DBModel {
	conversationId: number
	whatsappID: string
	direction: "in" | "out"
	footerMessage: string
	message: string
	headerMessage: string
	timestamp: number
	buttons: null | Array<MessageButton>
}

export interface MessageButton extends DBModel {
	text: string
	payload: string
}

interface ConversationsState {
	conversations: Array<Conversation>
	setConversations: (conversations: Array<Conversation>) => void
	newConversation: (conversation: Conversation) => void
	updateConversation: (conversation: Conversation) => void
	addMessage: (message: Message) => void
}

export const useConversationsStore = create<ConversationsState>((set) => ({
	conversations: [],
	setConversations(conversations) {
		set((state) => ({ ...state, conversations }))
	},
	newConversation(conversation) {
		set((state) => {
			for (let idx = 0; idx < state.conversations.length; idx++) {
				const conversationNeedle = state.conversations[idx]

				if (conversation.phoneNumberId === conversationNeedle.phoneNumberId) {
					const conversations = [...state.conversations]
					conversations[idx] = conversation
					return {
						...state,
						conversations,
					}
				}
			}

			return {
				...state,
				conversations: [...state.conversations, conversation],
			}
		})
	},
	updateConversation(conversation) {
		set((state) => {
			const conversations = [...state.conversations]
			for (let idx = 0; idx < conversations.length; idx++) {
				if (conversations[idx].ID === conversation.ID) {
					conversations[idx] = conversation
					break
				}
			}

			return {
				...state,
				conversations,
			}
		})
	},
	addMessage(message) {
		set((state) => {
			for (let idx = 0; idx < state.conversations.length; idx++) {
				const conversation = state.conversations[idx]
				if (message.conversationId === conversation.ID) {
					conversation.messages.push(message)
					state.conversations[idx] = conversation
					return {
						...state,
						conversations: [...state.conversations],
					}
				}
			}

			return state
		})
	},
}))
