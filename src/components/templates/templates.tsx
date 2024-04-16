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
import { useEffect, useState, Fragment } from "react"
import { DBModel, emptyDBModel } from "@/lib/types"
import { TrashIcon } from "@radix-ui/react-icons"
import { OpenCloseButton } from "../openCloseButton"

interface Template extends DBModel {
	name: string
	header: string | null
	body: string
	footer: string | null
	templateCustomButtons: Array<TemplateCustomButton>
}

interface TemplateCustomButton extends DBModel {
	templateID: number
	text: string
}

export function Templates() {
	const [open, setOpen] = useState(false)
	const [templates, setTemplates] = useState<Array<Template>>()
	const [newTemplateOpen, setNewTemplateOpen] = useState(false)

	const getData = async () => {
		const templatesResponse = await fetch("/api/templates")
		setTemplates(await templatesResponse.json())
	}

	// const updateTemplate = (index: number, template: Template) => {
	// 	setTemplates((prev) => {
	// 		if (!prev) return undefined
	// 		const templates = [...prev]
	// 		templates[index] = template
	// 		return templates
	// 	})
	// }

	const removeTemplate = async (index: number) => {
		if (!templates) return

		const template = templates[index]
		await fetch(`/api/templates/${template.ID}`, { method: "DELETE" })
		setTemplates((prev) => {
			if (!prev) return undefined

			const templates = [...prev]
			templates.splice(index, 1)
			return templates
		})
	}

	useEffect(() => {
		getData()
	}, [])

	const newTemplate = (template: Template) => {
		setTemplates((prev) => {
			if (!prev) return undefined

			return [...prev, template]
		})
	}

	return (
		<>
			<h2 m-6 mb-0 flex flex-wrap gap-4 justify-between items-center>
				<span inline-flex items-center>
					<OpenCloseButton open={open} setOpen={setOpen} /> Message templates
				</span>
				<Button onClick={() => setNewTemplateOpen(true)}>New template</Button>
			</h2>

			{templates && open ? (
				<div flex flex-wrap gap-4 p-4>
					{templates.map((template, idx) => (
						<Template
							template={template}
							key={template.ID}
							remove={() => removeTemplate(idx)}
						/>
					))}
				</div>
			) : undefined}
			<NewTemplateDialog
				open={newTemplateOpen}
				newTemplate={newTemplate}
				close={() => setNewTemplateOpen(false)}
			/>
		</>
	)
}

interface TemplateProps {
	template: Template
	remove: () => void
}

function Template({ template, remove }: TemplateProps) {
	return (
		<div flex w-full gap-4>
			<div>
				<Button onClick={remove} variant="ghost">
					<TrashIcon />
				</Button>
			</div>
			<div overflow-hidden>
				<h4>{template.name}</h4>
				<p truncate text-sm text-zinc-400>
					{template.body} {template.body}
				</p>
				<p text-sm text-zinc-500>
					{template.header ? (
						<span>
							Header: <span text-zinc-400>{template.header}</span>
						</span>
					) : undefined}{" "}
					{template.footer ? (
						<span>
							Footer: <span text-zinc-400>{template.footer}</span>
						</span>
					) : undefined}
					{template.templateCustomButtons.map((btn, idx) => (
						<Fragment key={idx}>
							{" "}
							<span>
								Button {idx + 1}: <span text-zinc-400>{btn.text}</span>
							</span>
						</Fragment>
					))}
				</p>
			</div>
		</div>
	)
}

interface NewTemplateDialogProps {
	open: boolean
	newTemplate: (template: Template) => void
	close: () => void
}

const emptyTemplate = (): Template => ({
	...emptyDBModel(),
	name: "hello_world_2",
	header: null,
	body: "",
	footer: null,
	templateCustomButtons: [],
})

function NewTemplateDialog({
	open,
	newTemplate,
	close,
}: NewTemplateDialogProps) {
	const [state, setState] = useState<Template>(emptyTemplate())

	const createConversation = async () => {
		const response = await post("/api/templates", state)
		const template = await response.json()
		newTemplate(template)
		setState(emptyTemplate())
	}

	return (
		<AlertDialog open={open} onOpenChange={() => close()}>
			<AlertDialogContent>
				<AlertDialogHeader>
					<AlertDialogTitle>Create a new conversation</AlertDialogTitle>
					<AlertDialogDescription></AlertDialogDescription>
				</AlertDialogHeader>
				<Label htmlFor="header">Name</Label>
				<Input
					value={state.name}
					onChange={(e) => setState((s) => ({ ...s, header: e.target.value }))}
					name="name"
					id="name"
					placeholder="hello_world"
				/>
				<Label htmlFor="header">Header</Label>
				<Input
					value={state.header ?? ""}
					onChange={(e) => setState((s) => ({ ...s, header: e.target.value }))}
					name="header"
					id="header"
					placeholder="Header"
				/>
				<Label htmlFor="body">Body</Label>
				<Input
					value={state.body}
					onChange={(e) => setState((s) => ({ ...s, body: e.target.value }))}
					name="body"
					id="body"
					placeholder="Hello world!"
				/>
				<Label htmlFor="footer">Footer</Label>
				<Input
					value={state.footer ?? ""}
					onChange={(e) => setState((s) => ({ ...s, footer: e.target.value }))}
					name="footer"
					id="footer"
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
