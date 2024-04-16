import {
	presetAttributify,
	presetWind,
	defineConfig,
	transformerDirectives,
} from "unocss"
import transformerAttributifyJsx from "@unocss/transformer-attributify-jsx-babel"

export default defineConfig({
	presets: [presetAttributify(), presetWind()],
	transformers: [transformerAttributifyJsx(), transformerDirectives()],
})
