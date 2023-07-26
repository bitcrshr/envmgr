<div class="w-full h-screen bg-mystic flex flex-col p-3 gap-4">
	<header class="w-full flex items-center rounded-md">
		<span class="font-barlow text-4xl text-pale-sky">envmgr</span>
	</header>

	<slot />
</div>

<style>
	.nm-container {
		border-radius: 20px;
		background: #d5e0e7;
		box-shadow: 20px 20px 60px #b5bec4, -20px -20px 60px #f5ffff;
	}
</style>
