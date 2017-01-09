import React, { Component } from 'react';
import request from "superagent"
const {ipcRenderer} = require('electron');

import FontIcon from 'material-ui/FontIcon';
import FloatingActionButton from 'material-ui/FloatingActionButton';

import NavBar from "./NavBar"
import Search from "./Search"
import ToolTip from 'react-portal-tooltip'
import {Table, TableBody, TableHeader, TableHeaderColumn, TableRow, TableRowColumn} from 'material-ui/Table';

import util from "../util"

class Downloads extends Component{

	constructor(props){
		super(props);

		this.state = {
			torrents: [],
			advanced: false
		}

		this.tick = this.tick.bind(this);
		this.componentDidMount = this.componentDidMount.bind(this);
		this.componentWillUnmount = this.componentWillUnmount.bind(this);
	}

	componentDidMount(){
		this.tick();
	}

	componentWillUnmount(){
		clearTimeout(this.timer);
	}

	tick(){
		hadouken.list((err, res) => {
			this.setState({ torrents: res });

			this.timer = setTimeout(this.tick, 3000);
		});
	}

	render() {
		return (
			<div style={{height: "100%"}}>
			<div className="parent">
				<NavBar advancedClick={()=>this.setState({ advanced: !this.state.advanced})} />
				{ !this.state.advanced &&
				<div>
					<Table fixedHeader={false} style={{tableLayout: "auto"}}>
						<TableHeader>
							<TableRow>
								<TableHeaderColumn>Name</TableHeaderColumn>
								<TableHeaderColumn>Progress</TableHeaderColumn>
								<TableHeaderColumn>Size</TableHeaderColumn>
								<TableHeaderColumn>Download</TableHeaderColumn>
								<TableHeaderColumn>Upload</TableHeaderColumn>
								<TableHeaderColumn>Seed/Peer</TableHeaderColumn>
							</TableRow>
						</TableHeader>

						<TableBody>
								{this.state.torrents.map((torrent, index) => {
									return (<TableRow key={index}>
												<TableRowColumn>{torrent.name ? torrent.name : torrent.infohash}</TableRowColumn>
												<TableRowColumn>{torrent.progress /10}%</TableRowColumn>
												<TableRowColumn>
													{util.bytes_to_size(torrent.size)}
												</TableRowColumn>
												<TableRowColumn>{util.bytes_to_size(torrent.down)}/s</TableRowColumn>
												<TableRowColumn>{util.bytes_to_size(torrent.up)}/s</TableRowColumn>
												<TableRowColumn>{torrent.listSeeds}/{torrent.listPeers}</TableRowColumn>
											</TableRow>)
								})}
						</TableBody>
					</Table>
				</div>
				}
			
				
				{ this.state.advanced &&
						<iframe 
							id="advancedTorrent"
							src="http://admin:admin@127.0.0.1:7070"/>
				}
				</div>
			</div>
		)
	}
}

export default Downloads;
