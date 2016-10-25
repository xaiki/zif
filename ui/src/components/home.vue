<template>
		<div class="contain">
			<div class="row">

				<div class="col s1"></div>

				<div class="col s10">
				<nav>
					<div class="nav-wrapper">
						<div class="input-field yellow darken-1 grey-text text-darken-2">
							<input id="search" type="search" required>
							 <label for="search"><i class="material-icons grey-text text-darken-2" >search</i></label>
							 <i class="material-icons">close</i>
						</div>
					</div>
				</nav>
				</div>

			</div>

			<div class="row">
				<div v-for="post in posts" class="col s12 m6 l4">
						<post   :title="post.Title">
						</post>
				</div>
			</div>
		</div>
</template>

<script>
import Post from "./post.vue"
import zif from "../zif.js"

export default{
	data() {
		return {
			posts: []
		}
	},

	methods: {
		refreshPosts: function() {
			zifd.recent(0, (data) => {
				this.posts = data.posts;
				console.log(data.posts)
			});
		}
	},

	created: function() {
		window.zifd = zif("127.0.0.1", "8080");
		this.refreshPosts();
	},

	components: {
		"post": Post
	},
}
</script>

<style>
body{
	background-color: #616161;
}
</style>
