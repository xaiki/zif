import React, { Component } from 'react';
import request from "superagent"
const {ipcRenderer} = require('electron');

import FontIcon from 'material-ui/FontIcon';

import NavBar from "./NavBar"
import Search from "./Search"
import ToolTip from 'react-portal-tooltip'
import {Table, TableBody, TableHeader, TableHeaderColumn, TableRow, TableRowColumn} from 'material-ui/Table';

import util from "../util"

class Downloads extends Component{

	constructor(props){
		super(props);

		this.state = {
			torrents: []
		}

		this.mounted = false;
		this.componentDidMount = this.componentDidMount.bind(this);
		this.componentWillUnmount = this.componentWillUnmount.bind(this);
		this.onTorrents = this.onTorrents.bind(this);
		this.onTorrent = this.onTorrent.bind(this);

		ipcRenderer.on("torrent", this.onTorrent);
		ipcRenderer.on("torrents", this.onTorrents);
	}

	onTorrent(e, torrent){
		if(this.mounted)
			this.setState({ torrents: this.state.torrents.push(torrent)});
	}

	onTorrents(e, torrents){
		if(this.mounted)
			this.setState({ torrents: torrents });
	}

	componentDidMount(){
		this.mounted = true;
		ipcRenderer.send("torrents");
	}

	componentWillUnmount(){
		this.mounted = false;

		ipcRenderer.removeListener("clienttorrent", this.onTorrent);
		ipcRenderer.removeListener("torrents", this.onTorrents);
	}

	render() {
		console.log(this.state.torrents)
		return (
			<div>
				<NavBar />
				<div style={{marginTop: "50px"}}>
					<Table fixedHeader={false} style={{tableLayout: "auto"}}>
						<TableHeader>
							<TableRow>
								<TableHeaderColumn>Name</TableHeaderColumn>
								<TableHeaderColumn>Progress</TableHeaderColumn>
								<TableHeaderColumn>Size</TableHeaderColumn>
								<TableHeaderColumn>Download</TableHeaderColumn>
								<TableHeaderColumn>Upload</TableHeaderColumn>
								<TableHeaderColumn>Seed/Leech</TableHeaderColumn>
							</TableRow>
						</TableHeader>

						<TableBody>
								{this.state.torrents.map((torrent, index) => {
									return (<TableRow key={index}>
												<TableRowColumn>{torrent.name ? torrent.name : torrent.infoHash}</TableRowColumn>
												<TableRowColumn>{torrent.progress}</TableRowColumn>
												<TableRowColumn>
													{torrent.info ? util.bytes_to_size(torrent.info.length) : <span>?</span>}
												</TableRowColumn>
												<TableRowColumn></TableRowColumn>
												<TableRowColumn></TableRowColumn>
												<TableRowColumn></TableRowColumn>
											</TableRow>)
								})}
						</TableBody>
					</Table>
				</div>
			</div>
		)
	}
}

export default Downloads;
