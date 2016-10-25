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
				<div class="col s12 m6 l4">
					<post title="Big Buck Bunny" description="A large and lovable rabbit deals with three tiny bullies, led by a flying squirrel, who are determined to squelch his happiness." img="https://upload.wikimedia.org/wikipedia/commons/c/c5/Big_buck_bunny_poster_big.jpg"></post>
				</div>
			</div>
		</div>
</template>

<script>
import Post from "./post.vue"
import zif from "../zif.js"

let zifd = zif("127.0.0.1", "8080");
console.log(zifd.search("Ze9fSR7QkZVWGddt2yPppJvx2ffxk2PUjW", "foo", (d)=>{console.log(d)}))

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
body{
	background-color: #616161;
}
</style>
