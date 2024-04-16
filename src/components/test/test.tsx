import { HelloWorldTemplateExample } from "./template"
import { HelloWorldMessage } from "./message"
import { State } from "@/services/state"
import { useState } from "react"
import { OpenCloseButton } from "../openCloseButton"

export interface TestProps {
	state: State
}

export function Test({ state }: TestProps) {
	const [open, setOpen] = useState(false)

	return (
		<>
			<h2 m-6 mb-0 flex flex-wrap gap-4 justify-between items-center>
				<span inline-flex items-center>
					<OpenCloseButton open={open} setOpen={setOpen} /> API Examples
				</span>
			</h2>

			{open ? (
				<div p-4 flex flex-col gap-6>
					<HelloWorldTemplateExample state={state} />
					<HelloWorldMessage state={state} />
				</div>
			) : undefined}
		</>
	)
}
