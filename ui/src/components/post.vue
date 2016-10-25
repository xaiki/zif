<template>
	<div class="card post hoverable">
		<div class="card-image">
			<img :src="img">
		</div>
		<div class="card-content">
			<h5>{{title}}</h5>
			<p>{{description}}</p>
		</div>
		<div class="card-action">
			<div class="row">
				<div class="col s6">
					<div class="waves-effect waves-light btn">Download</div>
				</div>
				<div class="col s6">
					<div class="waves-effect waves-light btn">Stream</div>
				</div>
			</div>
		</div>
	</div>
</template>	

<script>
export default{
	props: {
		title: String,
		description: String,
		img: String
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
