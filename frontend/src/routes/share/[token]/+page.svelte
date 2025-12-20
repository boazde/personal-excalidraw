<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { createQuery } from '@tanstack/svelte-query';
	import { drawingsAPI } from '$lib/api';
	import PublicExcalidrawWrapper from '$lib/components/PublicExcalidrawWrapper.svelte';
	import type { DrawingState } from '$lib/stores/drawing';

	const shareToken = page.params.token || '';

	// Query for public drawing (no auth required)
	const drawingQuery = createQuery(() => ({
		queryKey: ['public-drawing', shareToken],
		queryFn: () => drawingsAPI.getByShareToken(shareToken),
		enabled: !!shareToken,
		retry: false
	}));

	let initialData = $state<DrawingState | null>(null);

	onMount(async () => {
		if (!shareToken) {
			return;
		}

		// Load drawing content from share token
		try {
			const drawing = await drawingsAPI.getByShareToken(shareToken);
			initialData = {
				elements: (drawing.data?.elements || []) as unknown[],
				appState: (drawing.data?.appState || {}) as unknown,
				files: (drawing.data?.files || {}) as Record<string, unknown>
			};
		} catch (error) {
			console.error('Failed to load public drawing:', error);
		}
	});
</script>

<div class="h-screen w-screen flex flex-col bg-base-100">
	<!-- Header - matches drawing page exactly -->
	<div class="flex items-center justify-between px-4 py-2 border-b border-base-300">
		<!-- Read-only notice -->
		<div class="w-24 flex items-center">
			<div class="badge badge-ghost badge-sm gap-1">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="size-3"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
					/>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
					/>
				</svg>
				<span>View Only</span>
			</div>
		</div>

		<!-- Centered title -->
		<div class="flex-1 flex justify-center px-4">
			{#if drawingQuery.isFetching}
				<span class="loading loading-spinner loading-sm"></span>
			{:else if drawingQuery.data}
				<h1 class="text-lg font-semibold">{drawingQuery.data.name}</h1>
			{/if}
		</div>

		<!-- Right side status (matches drawing page) -->
		<div class="flex items-center gap-2">
			{#if drawingQuery.isFetching}
				<span class="loading loading-spinner loading-xs"></span>
				<span class="text-xs text-base-content/70">Loading...</span>
			{:else if drawingQuery.data}
				<span class="text-xs text-base-content/70">Ready</span>
			{/if}
		</div>
	</div>

	<!-- Drawing canvas -->
	<div class="flex-1 flex overflow-hidden">
		{#if drawingQuery.isError}
			<div class="flex flex-col items-center justify-center w-full h-full gap-4">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="size-16 text-error"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
					/>
				</svg>
				<div class="text-center">
					<h2 class="text-xl font-semibold mb-2">Drawing Not Found</h2>
					<p class="text-base-content/70">
						This shared drawing doesn't exist or is no longer available.
					</p>
				</div>
			</div>
		{:else if initialData}
			<PublicExcalidrawWrapper {initialData} />
		{:else}
			<div class="flex items-center justify-center w-full h-full">
				<span class="loading loading-spinner loading-lg"></span>
			</div>
		{/if}
	</div>
</div>
