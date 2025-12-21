<script lang="ts">
	import { goto } from '$app/navigation';
	import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
	import { drawingsAPI } from '$lib/api';
	import type { ID } from '$lib/types';
	import { toastStore } from '$lib/stores/toast';

	const queryClient = useQueryClient();

	// Query for drawings list
	const drawingsQuery = createQuery(() => ({
		queryKey: ['drawings'],
		queryFn: () => drawingsAPI.list({ limit: 100, offset: 0 })
	}));

	// Mutation for creating drawing
	const createDrawingMutation = createMutation(() => ({
		mutationFn: (name: string) =>
			drawingsAPI.create({
				name,
				data: { elements: [], appState: {}, files: {} }
			}),
		onSuccess: (data) => {
			queryClient.invalidateQueries({ queryKey: ['drawings'] });
			goto(`/drawing/${data.id}`);
		}
	}));

	// Mutation for updating drawing name
	const updateDrawingMutation = createMutation(() => ({
		mutationFn: ({ id, name }: { id: string; name: string }) => drawingsAPI.update(id, { name }),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ['drawings'] });
		}
	}));

	// Mutation for deleting drawing
	const deleteDrawingMutation = createMutation(() => ({
		mutationFn: (id: string) => drawingsAPI.delete(id),
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ['drawings'] });
		}
	}));

	// Mutation for duplicating drawing
	const duplicateDrawingMutation = createMutation(() => ({
		mutationFn: async (id: string) => {
			// Fetch original drawing
			const original = await drawingsAPI.get(id);

			// Create new drawing with copied data and modified name
			return drawingsAPI.create({
				name: `${original.name} - Copy`,
				data: original.data
			});
		},
		onSuccess: (data) => {
			queryClient.invalidateQueries({ queryKey: ['drawings'] });
			toastStore.show({
				type: 'success',
				message: `Drawing "${data.name}" created successfully`,
				duration: 3000
			});
		},
		onError: (error) => {
			toastStore.show({
				type: 'error',
				message: error instanceof Error ? error.message : 'Failed to duplicate drawing',
				duration: 5000
			});
		}
	}));

	// UI state (Svelte stores for UI only)
	let editingId = $state<ID | null>(null);
	let editingValue = $state('');

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function handleNewDrawing() {
		createDrawingMutation.mutate('Untitled Drawing');
	}

	function handleDelete(id: ID) {
		if (confirm('Delete this drawing?')) {
			deleteDrawingMutation.mutate(id);
		}
	}

	function handleDuplicate(id: ID) {
		duplicateDrawingMutation.mutate(id);
	}

	function startEditing(id: ID, currentName: string) {
		editingId = id;
		editingValue = currentName;
	}

	function saveEdit() {
		if (editingId === null || !editingValue.trim()) return;
		updateDrawingMutation.mutate({
			id: editingId,
			name: editingValue.trim()
		});
		editingId = null;
		editingValue = '';
	}

	function cancelEdit() {
		editingId = null;
		editingValue = '';
	}

	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			event.preventDefault();
			saveEdit();
		} else if (event.key === 'Escape') {
			cancelEdit();
		}
	}

	function copyShareURL(shareToken: string) {
		const shareURL = `${window.location.origin}/share/${shareToken}`;
		navigator.clipboard.writeText(shareURL);
		toastStore.show({
			type: 'success',
			message: 'Share link copied to clipboard!',
			duration: 3000
		});
	}
</script>

<div class="min-h-screen bg-base-100 p-8">
	<div class="max-w-6xl mx-auto">
		<div class="mb-8">
			<img src="/logo-text.svg" alt="Personal Excalidraw" class="h-6 w-auto" />
		</div>
		<div class="flex justify-between items-center mb-8">
			<h1 class="text-lg font-semibold">My Drawings</h1>
			<button
				onclick={handleNewDrawing}
				class="btn btn-sm btn-primary"
				disabled={createDrawingMutation.isPending}
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-5 w-5"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M12 4v16m8-8H4"
					/>
				</svg>
				New
			</button>
		</div>

		<!-- Loading State -->
		{#if drawingsQuery.isFetching}
			<div class="flex justify-center py-16">
				<span class="loading loading-spinner loading-lg"></span>
			</div>

			<!-- Error State -->
		{:else if drawingsQuery.isError}
			<div class="alert alert-error mb-4">
				<span>{drawingsQuery.error.message}</span>
				<button class="btn btn-sm" onclick={() => drawingsQuery.refetch()}> Retry </button>
			</div>

			<!-- Empty State -->
		{:else if drawingsQuery.data && drawingsQuery.data.drawings.length === 0}
			<div class="text-center py-16">
				<p class="text-lg">No drawings yet</p>
				<button onclick={handleNewDrawing} class="btn btn-primary btn-link mt-4">
					Create your first drawing
				</button>
			</div>

			<!-- Data Table -->
		{:else}
			<div class="overflow-x-auto">
				<table class="table w-full">
					<thead>
						<tr>
							<th>Name</th>
							<th>Created At</th>
							<th>Updated At</th>
							<th class="text-right">Actions</th>
						</tr>
					</thead>
					<tbody>
						{#if drawingsQuery.data}
							{#each drawingsQuery.data.drawings as drawing}
								<tr>
									<td class="font-medium">
										{#if editingId === drawing.id}
											<input
												type="text"
												class="input input-bordered input-sm w-full max-w-xs"
												bind:value={editingValue}
												onblur={saveEdit}
												onkeydown={handleKeydown}
											/>
										{:else}
											<button
												class="text-left hover:text-primary hover:underline cursor-pointer transition-colors"
												onclick={() => startEditing(drawing.id, drawing.name)}
												title="Click to edit name"
											>
												{drawing.name}
											</button>
										{/if}
									</td>
									<td>{formatDate(drawing.created_at)}</td>
									<td>{formatDate(drawing.updated_at)}</td>
									<td class="text-right">
										<div class="flex gap-2 justify-end">
											{#if drawing.share_token}
												<button
													class="btn btn-ghost btn-success btn-sm"
													title="Copy share link"
													onclick={() => copyShareURL(drawing.share_token!)}
													aria-label="Copy share link"
												>
													<svg
														xmlns="http://www.w3.org/2000/svg"
														class="size-4"
														fill="none"
														viewBox="0 0 24 24"
														stroke="currentColor"
													>
														<path
															stroke-linecap="round"
															stroke-linejoin="round"
															stroke-width="2"
															d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z"
														/>
													</svg>
												</button>
											{/if}
											<a
												href="/drawing/{drawing.id}"
												class="btn btn-ghost btn-info btn-sm"
												title="Edit drawing"
												aria-label="Edit drawing"
											>
												<svg
													xmlns="http://www.w3.org/2000/svg"
													class="size-4"
													fill="none"
													viewBox="0 0 24 24"
													stroke="currentColor"
												>
													<path
														stroke-linecap="round"
														stroke-linejoin="round"
														stroke-width="2"
														d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
													/>
												</svg>
											</a>
											<button
												class="btn btn-ghost btn-warning btn-sm"
												title="Duplicate drawing"
												onclick={() => handleDuplicate(drawing.id)}
												disabled={duplicateDrawingMutation.isPending}
												aria-label="Duplicate drawing"
											>
												<svg
													xmlns="http://www.w3.org/2000/svg"
													class="size-4"
													fill="none"
													viewBox="0 0 24 24"
													stroke="currentColor"
												>
													<path
														stroke-linecap="round"
														stroke-linejoin="round"
														stroke-width="2"
														d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
													/>
												</svg>
											</button>
											<button
												class="btn btn-ghost btn-error btn-sm"
												title="Delete drawing"
												onclick={() => handleDelete(drawing.id)}
												disabled={deleteDrawingMutation.isPending}
												aria-label="Delete drawing"
											>
												<svg
													xmlns="http://www.w3.org/2000/svg"
													class="size-4"
													fill="none"
													viewBox="0 0 24 24"
													stroke="currentColor"
												>
													<path
														stroke-linecap="round"
														stroke-linejoin="round"
														stroke-width="2"
														d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
													/>
												</svg>
											</button>
										</div>
									</td>
								</tr>
							{/each}
						{/if}
					</tbody>
				</table>
			</div>
		{/if}
	</div>
</div>
