export function CodeBlock({ code }: { code: string }) {
	return (
		<pre rounded bg-zinc-800 p-3 text-sm font-mono overflow-x-scroll>
			<br />
			{code}
			<br />
			<br />
		</pre>
	)
}
