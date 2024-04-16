import { toast } from "sonner"

export function getUrl(path: string, isWebsocket = false) {
	if (import.meta.env.DEV) {
		return `${isWebsocket ? "ws" : "http"}://localhost:1090${path}`
	}

	if (isWebsocket) {
		const protocol = location.protocol === "https:" ? "wss://" : "ws://"
		return protocol + location.host + path
	}

	return path
}

export async function fetch(path: string, options?: RequestInit) {
	const response = await window.fetch(getUrl(path), options)
	if (response.status >= 400) {
		let errMsg = await response.text()
		try {
			errMsg = JSON.parse(errMsg).error
		} catch (e) {
			// Ignore
		}
		toast.error(errMsg)
		throw new Error(errMsg)
	}
	return response
}

export function post(path: string, data: any, options?: RequestInit) {
	return fetch(path, {
		...options,
		method: "POST",
		headers: {
			"Content-Type": "application/json",
			...options?.headers,
		},
		body: JSON.stringify(data),
	})
}
