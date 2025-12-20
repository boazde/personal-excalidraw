<script lang="ts">
	import { onMount } from 'svelte'
	import type { DrawingState } from '$lib/stores/drawing'

	// Accept props from parent
	export let initialData: DrawingState

	let container: HTMLDivElement
	let excalidrawApp: any

	onMount(() => {
		// Dynamic import to avoid issues with React in Svelte
		const initExcalidraw = async () => {
			try {
				const { Excalidraw } = await import('@excalidraw/excalidraw')
				const React = await import('react')
				const ReactDOM = await import('react-dom/client')

				if (!container) return

				// Create a React root and render Excalidraw
				const root = ReactDOM.createRoot(container)

				// View-only mode: no onChange handler, no saving
				const ExcalidrawComponent = React.createElement(Excalidraw, {
					// No onChange - all edits are temporary (session-only)
					initialData: {
						elements: (initialData.elements || []) as any,
						appState: {
							...(initialData.appState || {}),
							collaborators: new Map(),
							// Enable view mode but allow temporary editing
							viewModeEnabled: false,
						},
						files: (initialData.files || {}) as any,
					},
					// UI options for public view
					UIOptions: {
						canvasActions: {
							// Disable loading
							loadScene: false,
						}
					}
				})

				root.render(ExcalidrawComponent)
				excalidrawApp = root
			} catch (error) {
				console.error('Failed to load Excalidraw:', error)
			}
		}

		initExcalidraw()

		// Return cleanup function for onMount
		return () => {
			if (excalidrawApp) {
				excalidrawApp.unmount()
			}
		}
	})
</script>

<div bind:this={container} class="w-full h-full"></div>

<style>
	:global(.excalidraw) {
		font-family: 'Cascadia Code', monospace;
	}

	/* Hide the library button */
	:global(.excalidraw .sidebar-trigger) {
		display: none !important;
	}
</style>
