import { ChevronDownIcon, ChevronUpIcon } from "@radix-ui/react-icons"
import { Button } from "./ui/button"

export interface OpenCloseButtonProps {
	open: boolean
	setOpen: (v: boolean) => void
}

export function OpenCloseButton({ open, setOpen }: OpenCloseButtonProps) {
	return (
		<Button size="sm" mr-2 variant="secondary" onClick={() => setOpen(!open)}>
			{open ? <ChevronUpIcon /> : <ChevronDownIcon />}
		</Button>
	)
}
