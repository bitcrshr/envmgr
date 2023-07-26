<!-- <script lang="ts">
	import { onMount } from 'svelte';
	import type { PageData } from './$types';

	export let data: PageData;

	onMount(() => {
		document.cookie = 'access_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
		document.cookie = 'refresh_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
		document.cookie = 'id_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';

		window.location.assign(data.logoutUrl);
	});
</script> -->
