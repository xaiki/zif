<template>
	<div class="contain">
		<h1>Stream</h1>
		<h1>{{$route.params.ih}}</h1>
	</div>
</template>

<script>
import util from "../util.js"

var WebTorrent = require('webtorrent')
var client = new WebTorrent()

export default{
	methods: {
		addMagnet: function(link) {
			console.log(link)

			client.add(link, function (torrent) {
			  // Got torrent metadata!
			  console.log('Client is downloading:', torrent.infoHash)

			  torrent.files.forEach(function (file) {
				// Display the file by appending it to the DOM. Supports video, audio, images, and
				// more. Specify a container element (CSS selector or reference to DOM node).
				file.appendTo('body')
			  })
			})
		}
	},
	
	created: function() {
		this.addMagnet(util.make_magnet(this.$route.params.ih))
	}
}
</script>
