<template>
	<div class="ui link cards">
		<post v-for="post in posts" v-bind:title="post.title"></post>
	</div>
</template>

<script>
import Post from "./post.vue"
import zif from "../zif.js"

let zifd = zif("127.0.0.1", "8080");

export default{
	data() {
		return {
		
		}
	},

	created: function() {
	},

	props: ["title"],

	components: {
		"post": Post
	}
}
</script>

<style>
body {
  font-family: Arial, sans-serif;
}
</style>
