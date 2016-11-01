<template>
	<div class="card post hoverable">
		<div class="card-image">
			<img :src="image">
		</div>
		<div class="card-content">
			<h5>{{util.trim(title, 72, true).trimmed}}</h5>

			<p class="grey-text lighten-2">
				{{util.trim(title, 72, true).left}}
			</p>	

			<span class="chip green-text"> {{seeders}}</span>
			<span class="chip red-text"> {{leechers}}</span>
			<a :href="util.make_magnet(infohash)">
				Magnet
			</a>

			<p>{{description}}</p>
		</div>
		<div class="card-action">
			<router-link :to="{ name: 'stream', params: { ih: infohash }}">Stream</router-link>
		</div>
	</div>
</template>	

<script>
import util from "../util.js"

export default{
	data() {
		return {
			util: util,
			description: "No description"
		}
	},
	props: {
		title: String,
		infohash: String,
		description: String,
		image: String,
		seeders: Number,
		leechers: Number 
	}
}
</script>

<style scoped>
.post {
	padding-left: 10px;
	padding-right: 10px;
	height: 100%;
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
