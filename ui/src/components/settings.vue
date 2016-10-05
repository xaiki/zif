<template>
	<div class="settings">
		<h1>Settings</h1>
		<hr>

		<div class="bootstrap">
			<h3>Bootstrap</h3>
			<div class="ui raised segment">
				<p>Bootstrap your routing table from that of another peer, enter the
				peer's domain name or IP address below.</p>

				<div class="ui fluid action input">
					<input v-model="address" id="bootstrap_addr" type="text" placeholder="example.com">
					<button v-on:click="bootstrap()" class="ui button">Bootstrap</button>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
import zif from "../zif.js"

export default {
	data () {
		return { address: "foo" }
	},
	methods: {
		bootstrap: () => { console.log(this.default.data().address) }
	}
}

</script>

<style>
</style>
