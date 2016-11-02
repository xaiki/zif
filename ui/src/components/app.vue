<template>
	<div class="pushable" id="app">
		<div class="row">
			<div class="col s3">
				<ul class="side-nav fixed ">
					<li>
						<router-link to="/">
							<i class="material-icons">home</i> Home
						</router-link>
					</li>
				</ul>
			</div>

			<div class="col s9">
				<router-view class="pusher view"></router-view>
			</div>
		</div>

		<div class="fixed-action-btn" style="right: 24px;">
			<router-link :to="'/'" class="btn-floating btn-large amber lighten-1">
		  		<i class="large grey-text text-darken-2">Z</i>
			</router-link>
		</div>
	</div>
</template>

<script>
	window.zifd_address = "http://127.0.0.1:8080/";

	export default {
	}

</script>

<style lang="sass">
.row .side-nav {
	width: 250px;
}
body{
	background-color: #616161;
}
</style>
