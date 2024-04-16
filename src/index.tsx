/* @refresh reload */
import "./theme.css"
import "./global.css"

import "virtual:uno.css"

import { App } from "./pages/App"
import React from "react"
import ReactDOM from "react-dom/client"
import { Toaster } from "@/components/ui/sonner"

ReactDOM.createRoot(document.getElementById("root")!).render(
	<React.StrictMode>
		<App />
		<Toaster />
	</React.StrictMode>,
)
