<template>
	<div class="pushable" id="app">
		<router-view class="pusher view"></router-view>

		<div class="fixed-action-btn" style="right: 24px;">
			<a class="btn-floating btn-large amber lighten-1">
		  		<i class="large grey-text text-darken-2">Z</i>
			</a>
		</div>
	</div>
</template>

<script>
	import Sidebar from "./sidebar.vue"

	window.zifd_address = "http://127.0.0.1:8080/";

	export default {
		components: {
			"sidebar": Sidebar
		}
	}

</script>

<style lang="sass">
body{
	background-color: #616161;
}
</style>
