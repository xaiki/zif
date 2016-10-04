<template>
	<div class="settings">
		<h1>Settings</h1>
		<hr>

		<div class="bootstrap">
			<h3>Bootstrap</h3>
			<div class="ui raised segment">
				<p>Bootstrap your routing table from that of another peer, enter the
				peer's domain name or IP address below.</p>

				<div v-bind:class="{ loading: isLoading }" class="ui fluid action input">
					<input type="text" placeholder="example.com">
					<button class="ui button">Bootstrap</button>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
	import zif from "../zif.js"

	export default {
		data() {
			return {
				isLoading: true
			}
		}
	}
</script>

<style>

</style>
