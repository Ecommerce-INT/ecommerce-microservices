import { paraglideVitePlugin } from '@inlang/paraglide-js';
import adapter from '@sveltejs/adapter-bun';
import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	server: { port: 3000 },
	preview: { port: 3000 },
	plugins: [
		tailwindcss(),
		sveltekit({
			// Kit 3 replaced the built-in `$lib` alias with the `#lib` imports map.
			// We keep `$lib` because shadcn-svelte generates `$lib/...` imports.
			alias: { $lib: 'src/lib' },
			experimental: {
				remoteFunctions: true
			},
			compilerOptions: {
				experimental: {
					async: true
				},
				// Force runes mode for the project, except for libraries. Can be removed in svelte 6.
				runes: ({ filename }) =>
					filename.split(/[/\\]/).includes('node_modules') ? undefined : true
			},
			adapter: adapter()
		}),
		paraglideVitePlugin({
			project: './project.inlang',
			outdir: './src/lib/paraglide',
			// URL routes (/ for vi, /en/... for en) win, then the locale cookie.
			strategy: ['url', 'cookie', 'baseLocale'],
			emitTsDeclarations: true
		})
	]
});
