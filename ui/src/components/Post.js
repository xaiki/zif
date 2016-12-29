import React, { Component } from 'react';
import Chip from 'material-ui/Chip';
import { hashHistory, Link } from 'react-router';
import {Card, CardActions, CardHeader, CardText} from 'material-ui/Card';
import FlatButton from 'material-ui/FlatButton';
import {List, ListItem} from 'material-ui/List';

import moment from "moment"

const style = {
	marginTop: "10px"
};

import util from "../util"

class Post extends Component
{
	constructor(props){
		super(props);

		this.onContextMenu= this.onContextMenu.bind(this);
	}

	static get defaultProps() {
		return {
			title: "Untitled",
			infohash: "nohash",
			seeders: 0,
			leechers: 0,
			uploadDate: 0,
			fileCount: 0,
			size: 0,
			tags: [],
			source: "foo"
		}
	}

	formatUnixTime(time) {
		var m = moment.unix(time);

		return m.fromNow();
	}

	onContextMenu(e) {
		e.preventDefault();

		var menu = new Menu();	
		menu.append(new MenuItem({ 
			label: "Copy Magnet Link",
			click: () => {
				clipboard.writeText(util.make_magnet(this.props.infohash));
			}
		}));
		menu.popup(remote.getCurrentWindow());
	}

	render() {
		return (
			<div className="card" onContextMenu={this.onContextMenu}>
				<details>
					<summary className="header">
						<div style={{display: "inline"}}>
							<h2 className="title">{this.props.title}</h2>
							<div className="info">
								<div><span style={{color:"#279c10"}}>{this.props.seeders}</span> / <span style={{color:"#b11106"}}>{this.props.leechers}</span></div>
								<div>{this.props.fileCount} files, {util.bytes_to_size(this.props.size)}</div>
								<div>{this.formatUnixTime(this.props.uploadDate)}</div>
							</div>
						</div>
						<div className="source"><em>uploaded by {this.props.source}</em></div>
					</summary>

					<div className="body">
						<div className="description">
							<em>{this.props.meta.description == undefined &&
								"No description"}</em>
						</div>

						<div className="info">
							<a className="magnet"
								onContextMenu={this.onContextMenu}
								href={util.make_magnet(this.props.infohash)}>
								<i className="material-icons">link</i>
								<span> Magnet</span>
							</a>
						</div>
					</div>

				</details>
			</div>)
	}
}

export default Post;
