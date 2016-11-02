<template>
	<div class="container">
		<ul class="collection with-header">
			<li class="collection-header">
				<h5>{{$route.params.title}}</h5>
				{{util.bytes_to_size($route.params.size)}}, {{$route.params.filecount}} file/s
			</li>

			<li v-for="e in entries" class="collection-item">
				<entry :filename="e.name"></entry>
			</li>

		</ul>
	</div>
</template>

<script>
import entry from "./stream-entry.vue"
import util from "../util.js"
import stream from "../stream.js"

export default{
	data() {
		return {
			streamer: stream(this),
			util: util,
			entries: []
		}
	},
	methods: {
		addMagnet: function(link) {
			console.log(link)
			this.streamer.add(link)

			  /*console.log('Client is downloading:', torrent.infoHash)

			  torrent.files.forEach(function (file) {
			  	  var elem = document.createElement("li");
				  elem.className += "collection-item";
				  elem.innerHTML = file.name;

				  file.appendTo(elem, { autoplay: false }, (err, e) =>{
					$("#contents").append(elem);
					e.className += "secondary-content";
				  })
			  })
			})*/
		}
	},
	
	created: function() {
		this.addMagnet(util.make_magnet(this.$route.params.ih))
	},
	components: {
		"entry": entry
	}
}
</script>

<style>
</style>
