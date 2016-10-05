<template>
	<div class="pushable" id="app">
		<sidebar></sidebar>
		<router-view class="pusher view"></router-view>
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

<style>
.view {
	margin-left: 5%;
	width: 80%;
}
</style>
