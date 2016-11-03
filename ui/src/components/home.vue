<template>
		<div class="container">
			<div class="row">

				<div class="col s1"></div>

				<div class="col s10">
				<nav>
					<div class="nav-wrapper">
						<div class="input-field yellow darken-2 grey-text text-darken-2">
							<input id="search" type="search" v-on:keydown.enter="search" required>
							 <label for="search"><i class="material-icons grey-text text-darken-2" >search</i></label>
							 <i class="material-icons">close</i>
						</div>
					</div>
				</nav>
				</div>

			</div>

			<div v-for="post in posts" class="collection">

					<post   class="collection-item"
							:title="post.Title"
							:infohash="post.InfoHash"
							:seeders="post.Seeders"
							:leechers="post.Leechers"
							:size="post.Size"
							:filecount="post.FileCount"
					</post>

			</div>

		</div>
</template>

<script>
import Post from "./post.vue"
import zif from "../zif.js"
import util from "../util.js"

export default{
	data() {
		return {
			posts: [],
			util: util
		}
	},

	methods: {
		refreshPosts: function() {
			zifd.popular(0, (data) => {
				this.posts = data.posts;
				console.log(data.posts)
			});
		},
		search: function() {
			var query = $("#search").val()
			console.log("Query for: " + query)

			zifd.search(query, 0, (data) => {
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
</style>
