<template>
	<div class="settings ui divided items">
		<h1 class="header">Settings</h1>
		<div class="item bootstrap">

			<div class="content">
				<div class="header">Bootstrap</div>
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

let zifd = zif("127.0.0.1", "8080");

export default {
	data () {
		return { address: "foo" }
	},
	methods: {
		bootstrap: function() {
			$(".bootstrap button").addClass("loading");
			zifd.bootstrap(this.address, (data) => {
				$(".bootstrap button").removeClass("loading");
			});
		}
	}
}

</script>

<style>
</style>
