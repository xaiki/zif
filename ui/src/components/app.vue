<template>
	<div class="pushable" id="app">
		<sidebar></sidebar>
		<router-view class="pusher view"></router-view>
	</div>
</template>

<script>
	import Sidebar from "./sidebar.vue"

	export default {
		components: {
			"sidebar": Sidebar
		}
	}

</script>

<style>
body {
	font-family: Arial, sans-serif;
}
.view {
	margin-left: 5%;
	width: 80%;
}
</style>
