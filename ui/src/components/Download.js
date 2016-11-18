import React, { Component } from 'react';
import File from "./File";
import util from "../util"

var WebTorrent = require("webtorrent");

var client;

class Stream extends Component
{
	constructor(props) 
	{
		super(props);

		this.infohash = this.props.routeParams.infohash;
		this.magnet = this.infohashToMagnet(this.infohash);

		this.state = { files: [] };

		this.torrentAdded = this.torrentAdded.bind(this);
	}

	componentDidMount()
	{
		window.downloadClient.forEach((torrent) => { console.log(torrent);})
	}

	torrentAdded(torrent)
	{
		this.torrent = torrent;

		this.setState({ files: torrent.files.sort(
			(a, b) => {
				return util.sort.alphanum(a.name, b.name);
			}) 
		});
	}

	infohashToMagnet(ih)
	{
		return "magnet:?xt=urn:btih:" + ih;
	}

	render() {
		var i= 0;
		return(
			<div>
				{this.state.files.map((file) => {
					return <File 	key={i++}
									Index={i}
									Torrent={this.torrent}
									File={file}
									Title={file.name}/>
				})}
			</div>
		)
	}
}

export default Stream;
