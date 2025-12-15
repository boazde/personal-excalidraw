<script lang="ts">
	import { toastStore, type Toast } from '$lib/stores/toast';

	let toasts = $state<Toast[]>([]);

	$effect(() => {
		const unsubscribe = toastStore.subscribe((value) => {
			toasts = value;
		});
		return unsubscribe;
	});

	function getAlertClass(type: Toast['type']): string {
		const baseClass = 'alert alert-soft';
		switch (type) {
			case 'success':
				return `${baseClass} alert-success`;
			case 'error':
				return `${baseClass} alert-error`;
			case 'warning':
				return `${baseClass} alert-warning`;
			case 'info':
				return `${baseClass} alert-info`;
			default:
				return baseClass;
		}
	}
</script>

<!-- Toast Container -->
<div class="toast toast-bottom toast-end z-50">
	{#each toasts as toast (toast.id)}
		<div class={getAlertClass(toast.type)}>
			<span>{toast.message}</span>
			<button
				class="btn btn-sm btn-ghost btn-circle"
				onclick={() => toastStore.dismiss(toast.id)}
				aria-label="Dismiss notification"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-4 w-4"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M6 18L18 6M6 6l12 12"
					/>
				</svg>
			</button>
		</div>
	{/each}
</div>
