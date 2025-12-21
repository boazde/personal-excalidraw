<script lang="ts">
	import { createMutation, useQueryClient } from '@tanstack/svelte-query'
	import { drawingsAPI } from '$lib/api'
	import type { ID } from '$lib/types'

	let { open = $bindable(false), shareToken, drawingId } = $props<{
		open: boolean
		shareToken?: string
		drawingId: ID
	}>()

	const queryClient = useQueryClient()

	let modal = $state<HTMLDialogElement>()
	let copySuccess = $state(false)

	// Mutation for generating share token
	const generateTokenMutation = createMutation(() => ({
		mutationFn: () => drawingsAPI.generateShareToken(drawingId),
		onSuccess: () => {
			// Invalidate the drawing query to refetch with the new share token
			queryClient.invalidateQueries({ queryKey: ['drawing', drawingId] })
		}
	}))

	// Mutation for revoking share token
	const revokeTokenMutation = createMutation(() => ({
		mutationFn: () => drawingsAPI.revokeShareToken(drawingId),
		onSuccess: () => {
			// Invalidate the drawing query to refetch without the share token
			queryClient.invalidateQueries({ queryKey: ['drawing', drawingId] })
		}
	}))

	// Construct the shareable URL
	const shareUrl = $derived(
		shareToken ? `${typeof window !== 'undefined' ? window.location.origin : ''}/share/${shareToken}` : ''
	)

	// Open/close modal based on open prop
	$effect(() => {
		if (modal) {
			if (open) {
				modal.showModal()
			} else {
				modal.close()
			}
		}
	})

	async function copyToClipboard() {
		if (!shareUrl) return

		try {
			await navigator.clipboard.writeText(shareUrl)
			copySuccess = true
			setTimeout(() => {
				copySuccess = false
			}, 3000)
		} catch (err) {
			// Fallback: select the text
			const input = document.querySelector('#share-url-input') as HTMLInputElement
			if (input) {
				input.select()
			}
		}
	}

	function handleClose() {
		open = false
		copySuccess = false
	}

	function handleRevoke() {
		if (confirm('Are you sure you want to revoke this share URL? Anyone with the old link will no longer be able to access this drawing.')) {
			revokeTokenMutation.mutate()
		}
	}
</script>

<dialog bind:this={modal} class="modal" onclose={handleClose}>
	<div class="modal-box">
		<h3 class="text-lg font-bold mb-4">Share Drawing</h3>

		<p class="mb-4 text-sm text-base-content/70">
			Generate a shareable link that allows anyone with the link to view your drawing.
		</p>

		<!-- Share Link Input Section -->
		<div class="mb-6">
			{#if shareToken}
				<div class="join w-full">
					<div
						id="share-url-input"
						class="input input-bordered join-item font-mono text-ellipsis overflow-hidden pr-8 text-sm flex-1"
					>
						<span>
							{shareUrl}
						</span>
					</div>
					<button onclick={handleRevoke} class="btn btn-error join-item" title="Revoke share URL">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="1.5"
							stroke="currentColor"
							class="w-5 h-5"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"
							/>
						</svg>
						Revoke
					</button>
				</div>
			{:else}
				<div class="w-full">
					<div class="input input-bordered w-full text-sm text-base-content/50">
						Click 'Generate Share URL' button below to create a shareable link
					</div>
				</div>
			{/if}
		</div>

		{#if generateTokenMutation.isError}
			<div class="alert alert-error mb-4">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="1.5"
					stroke="currentColor"
					class="w-5 h-5"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z"
					/>
				</svg>
				<span>Failed to generate share link. Please try again.</span>
			</div>
		{/if}

		{#if revokeTokenMutation.isError}
			<div class="alert alert-error mb-4">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="1.5"
					stroke="currentColor"
					class="w-5 h-5"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z"
					/>
				</svg>
				<span>Failed to revoke share link. Please try again.</span>
			</div>
		{/if}

		{#if shareToken}
			<p class="alert alert-soft alert-warning flex items-start gap-2">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="1.5"
					stroke="currentColor"
					class="w-5 h-5 shrink-0"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z"
					/>
				</svg>
				<span
					>Revoking will invalidate the current share link. Anyone with the old link will no longer
					be able to access this drawing.</span
				>
			</p>
		{/if}

		<div class="modal-action">
			{#if !shareToken}
				<button
					onclick={() => generateTokenMutation.mutate()}
					disabled={generateTokenMutation.isPending}
					class="btn btn-primary"
				>
					{#if generateTokenMutation.isPending}
						<span class="loading loading-spinner loading-sm"></span>
						Generating...
					{:else}
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="1.5"
							stroke="currentColor"
							class="w-5 h-5"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								d="M7.217 10.907a2.25 2.25 0 100 2.186m0-2.186c.18.324.283.696.283 1.093s-.103.77-.283 1.093m0-2.186l9.566-5.314m-9.566 7.5l9.566 5.314m0 0a2.25 2.25 0 103.935 2.186 2.25 2.25 0 00-3.935-2.186zm0-12.814a2.25 2.25 0 103.933-2.185 2.25 2.25 0 00-3.933 2.185z"
							/>
						</svg>
						Generate Share URL
					{/if}
				</button>
			{:else}
				<button
					onclick={copyToClipboard}
					class="btn btn-primary"
					title={copySuccess ? 'Copied!' : 'Copy share URL to clipboard'}
				>
					{#if copySuccess}
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="1.5"
							stroke="currentColor"
							class="w-5 h-5 text-success"
						>
							<path stroke-linecap="round" stroke-linejoin="round" d="M4.5 12.75l6 6 9-13.5" />
						</svg>
						Copied!
					{:else}
						<svg
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
							stroke-width="1.5"
							stroke="currentColor"
							class="w-5 h-5"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								d="M15.666 3.888A2.25 2.25 0 0013.5 2.25h-3c-1.03 0-1.9.693-2.166 1.638m7.332 0c.055.194.084.4.084.612v0a.75.75 0 01-.75.75H9a.75.75 0 01-.75-.75v0c0-.212.03-.418.084-.612m7.332 0c.646.049 1.288.11 1.927.184 1.1.128 1.907 1.077 1.907 2.185V19.5a2.25 2.25 0 01-2.25 2.25H6.75A2.25 2.25 0 014.5 19.5V6.257c0-1.108.806-2.057 1.907-2.185a48.208 48.208 0 011.927-.184"
							/>
						</svg>
						Copy Share URL
					{/if}
				</button>
			{/if}
			<button onclick={handleClose} class="btn">Close</button>
		</div>
	</div>
	<form method="dialog" class="modal-backdrop">
		<button onclick={handleClose}>close</button>
	</form>
</dialog>
