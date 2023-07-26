<script lang="ts">
	import { goto } from '$app/navigation';

	let loginButtonPressed = false;
</script>

<div class="w-full h-screen flex items-center justify-center bg-mystic">
	<div
		class="nm-container w-1/2 max-w-[500px] min-w-[300px] h-[500px] flex flex-col items-center justify-center gap-20 p-8"
	>
		<span class="w-full flex items-center justify-center barlow text-5xl text-pale-sky select-none"
			>envmgr</span
		>

		<div
			role="button"
			tabindex="0"
			class="nm-button w-[200px] h-[60px] transition-all duration-250 flex items-center justify-center text-mystic text-xl select-none"
			class:pressed={loginButtonPressed}
			on:mousedown={() => (loginButtonPressed = true)}
			on:mouseup={() => (loginButtonPressed = false)}
			on:click={() => goto('/login')}
		>
			login
		</div>
	</div>
</div>

<style>
	.nm-container {
		border-radius: 50px;
		background: #d5e0e7;
		box-shadow: 20px 20px 60px #b5bec4, -20px -20px 60px #f5ffff;
	}

	.nm-button {
		border-radius: 50px;
		background: linear-gradient(145deg, #94a6b4, #7c8c97);
		box-shadow: 20px 20px 60px #75848f, -20px -20px 60px #9fb2c1;
	}

	.nm-button.pressed {
		border-radius: 50px;
		background: linear-gradient(145deg, #7c8c97, #94a6b4);
		box-shadow: 20px 20px 60px #75848f, -20px -20px 60px #9fb2c1;
	}
</style>
