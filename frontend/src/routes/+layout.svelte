<script lang="ts">
	import './layout.css'
	import favicon from '$lib/assets/favicon.svg'
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query'
	import { browser } from '$app/environment'
	import { page } from '$app/stores'
	import { onMount } from 'svelte'
	import { authStore } from '$lib/stores/auth'
	import AuthModal from '$lib/components/AuthModal.svelte'
	import Toast from '$lib/components/Toast.svelte'

	let { children } = $props()

	let showAuthModal = $state(false)
	let authChecked = $state(false)
	let isPublicRoute = $state(false)

	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				enabled: browser,
				staleTime: 1000 * 60,
				refetchOnWindowFocus: false,
				retry: (failureCount, error) => {
					// Don't retry on auth errors
					if (error instanceof Error && error.message.includes('Unauthorized')) {
						return false
					}
					return failureCount < 3
				}
			}
		}
	})

	// Check auth only in browser, after mount (prevents SSR hydration issues)
	onMount(() => {
		// Check if current route is public (share or public routes)
		const unsubscribe = page.subscribe(($page) => {
			isPublicRoute = $page.url.pathname.startsWith('/share/') || $page.url.pathname.startsWith('/public/')

			// Skip auth check for public routes
			if (isPublicRoute) {
				showAuthModal = false
				authChecked = true
			} else {
				const hasKey = authStore.hasAccessKey()
				showAuthModal = !hasKey
				authChecked = true
			}
		})

		// Listen for auth errors from API calls (but ignore on public routes)
		const handleAuthRequired = () => {
			if (!isPublicRoute) {
				showAuthModal = true
			}
		}
		window.addEventListener('auth:required', handleAuthRequired)

		return () => {
			unsubscribe()
			window.removeEventListener('auth:required', handleAuthRequired)
		}
	})

	function handleAuthSuccess() {
		showAuthModal = false
	}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{#if !browser || !authChecked}
	<!-- SSR or initial load: show loading -->
	<div class="flex items-center justify-center min-h-screen">
		<span class="loading loading-spinner loading-lg"></span>
	</div>
{:else if showAuthModal}
	<!-- Not authenticated: show modal only -->
	<AuthModal onSuccess={handleAuthSuccess} />
{:else}
	<!-- Authenticated: render app -->
	<QueryClientProvider client={queryClient}>
		{@render children()}
	</QueryClientProvider>
{/if}

<Toast />
