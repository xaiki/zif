import React, { Component } from 'react';

var WebTorrent = require("webtorrent");
var client;

class Stream extends Component
{
	constructor(props) 
	{
		super(props);

		this.infohash = this.props.routeParams.infohash;
		this.magnet = this.infohashToMagnet(this.infohash);

		this.torrentAdded = this.torrentAdded.bind(this);
	}

	componentDidMount()
	{
		client = new WebTorrent();
		console.log(client.add("magnet:?xt=urn:btih:d169e8930f820496bef1f42097f3ebcf6192fe52",
					function(t) {
						console.log(t);
					}));
	}

	torrentAdded(torrent)
	{
		console.log("added");
		console.log(torrent);
	}

	infohashToMagnet(ih)
	{
		return "magnet:?xt=urn:btih:" + ih;
	}

	render() {
		return(
				<h1></h1>
		)
	}
}

export default Stream;
