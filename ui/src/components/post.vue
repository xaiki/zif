<template>
	<div class="card post hoverable">
		<div class="card-image">
			<img :src="image">
		</div>
		<div class="card-content">
			<h5>{{title}}</h5>
			<span class="chip green-text"> {{seeders}}</span>
			<span class="chip red-text"> {{leechers}}</span>
			<a :href="util.make_magnet(infohash)"><i class="chip material-icons">link</i></a>
					
			<p>{{description}}</p>
		</div>
		<div class="card-action">
			<div class="row">
				<div class="col s4">
				</div>
			</div>
			<div class="row">
			</div>
		</div>
	</div>
</template>	

<script>
import util from "../util.js"

export default{
	data() {
		return {
			util: util
		}
	},
	props: {
		title: String,
		infohash: String,
		description: String,
		image: String,
		seeders: String,
		leechers: String
	}
}
</script>

<style scoped>
.post {
	padding-left: 10px;
	padding-right: 10px;
}

.card-image img {
	max-height: 300px;
	width: auto !important;
	margin: 0 auto;
}

.btn {
	margin-left: 3px;
	margin-right: 3px;
}
</style>
