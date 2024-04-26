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
import { Textarea } from "@/components/ui/textarea"

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

	const addButton = () =>
		setState((s) => {
			let text = "Button text"
			if (s.templateCustomButtons.length > 0) {
				text += " " + (s.templateCustomButtons.length + 1)
			}

			s.templateCustomButtons.push({
				...emptyDBModel(),
				templateID: s.ID,
				text,
			})

			return { ...s }
		})

	const setButtonText = (idx: number, text: string) =>
		setState((s) => {
			s.templateCustomButtons[idx].text = text
			return { ...s }
		})

	const removeButton = (idx: number) =>
		setState((s) => {
			s.templateCustomButtons.splice(idx, 1)
			return { ...s }
		})

	const intermediateClose = () => {
		close()
		setState(emptyTemplate())
	}

	return (
		<AlertDialog open={open} onOpenChange={() => intermediateClose()}>
			<AlertDialogContent>
				<AlertDialogHeader>
					<AlertDialogTitle>Create a new conversation</AlertDialogTitle>
				</AlertDialogHeader>
				<Label htmlFor="header">Name</Label>
				<Input
					value={state.name}
					onChange={(e) => setState((s) => ({ ...s, name: e.target.value }))}
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
				<Textarea
					value={state.body}
					onChange={(e) => setState((s) => ({ ...s, body: e.target.value }))}
					name="body"
					id="body"
					placeholder="Hello world!"
					h-30
				/>
				<Label htmlFor="footer">Footer</Label>
				<Input
					value={state.footer ?? ""}
					onChange={(e) => setState((s) => ({ ...s, footer: e.target.value }))}
					name="footer"
					id="footer"
					placeholder="Hello world!"
				/>

				{state.templateCustomButtons.length ? (
					<Label htmlFor="footer">Buttons</Label>
				) : undefined}
				{state.templateCustomButtons.map((btn, idx) => (
					<div key={idx} flex w-full items-center gap-4>
						<div>
							<Button onClick={() => removeButton(idx)} variant="ghost">
								<TrashIcon />
							</Button>
						</div>
						<div flex-1>
							<Label htmlFor="footer">Button #{idx + 1}</Label>
							<Input
								value={btn.text}
								onChange={(e) => setButtonText(idx, e.target.value)}
								name={"button-" + idx}
								id={"button-" + idx}
								placeholder="Hello world!"
							/>
						</div>
					</div>
				))}
				<div>
					<Button variant="secondary" onClick={addButton}>
						New button
					</Button>
				</div>
				<AlertDialogFooter>
					<AlertDialogCancel>Cancel</AlertDialogCancel>
					<AlertDialogAction onClick={createConversation}>
						Create
					</AlertDialogAction>
				</AlertDialogFooter>
			</AlertDialogContent>
		</AlertDialog>
	)
}
