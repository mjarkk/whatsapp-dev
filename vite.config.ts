import { defineConfig } from "vite"
import unocss from "unocss/vite"
import Icons from "unplugin-icons/vite"
import react from "@vitejs/plugin-react"
import path from "path"

export default defineConfig({
	plugins: [Icons({ compiler: "jsx", jsx: "react" }), unocss(), react()],
	resolve: {
		alias: {
			"@": path.resolve(__dirname, "./src"),
		},
	},
	server: {
		port: 3001,
	},
})
