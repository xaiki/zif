import React, { Component } from 'react';
import { Router, Route, hashHistory, Link } from 'react-router';
import request from "superagent"

import {List, ListItem} from 'material-ui/List';
import Divider from 'material-ui/Divider';
import Subheader from 'material-ui/Subheader';
import FontIcon from 'material-ui/FontIcon';
import IconButton from 'material-ui/IconButton';

import util from "../util.js"
import Upload from "./UploadDialog.js"
import EditAccount from "./EditAccount.js"
import ToolTip from 'react-portal-tooltip'

class NavBar extends Component{

	constructor(props){
		super(props);


		this.state = { name: window.entry.name, uploadTooltip: false, channelTooltip: false };
		window.navbar = this;
		this.uploadDialog = <Upload />;

		this.upload = {};
	}

	toggleUploadTooltip(){
		this.setState({ uploadTooltip: !this.state.uploadTooltip, channelTooltip: false })
	}

	toggleChannelTooltip(){
		request.get("http://127.0.0.1:8080/self/get/postcount/")
				.end(((err, res) => {
					if (err || res.body.status != "ok")
						return;

					window.entry.postCount = res.body.value;
					this.setState({});
				}).bind(this));

		this.setState({ channelTooltip: !this.state.channelTooltip, uploadTooltip: false })
	}

	render(){
		return(
				<ul className="topnav" id="mainMenu">
					<li><span id="logo">Zif</span></li>

					<li><Link to={"/"}>Home</Link></li>

					<li style={{float: "right"}}>
						<a id="channel"
							onMouseEnter={this.toggleChannelTooltip.bind(this)}
							onMouseLeave={this.toggleChannelTooltip.bind(this)}>
							{this.state.name}
						</a>
					</li>

					<li style={{float: "right", height:0}}>
						<a id="upload" 
								onMouseEnter={this.toggleUploadTooltip.bind(this)}
								onMouseLeave={this.toggleUploadTooltip.bind(this)}
								onClick={this.upload.open}>
							<i className="material-icons">file_upload</i>
						</a>
					</li>


					<ToolTip active={this.state.uploadTooltip} 
							position="bottom" 
							arrow="center" 
							parent="#upload">
						<div>
							Upload
						</div>
					</ToolTip>

					<ToolTip active={this.state.channelTooltip} 
							position="bottom" 
							arrow="center" 
							parent="#channel">
						<div>
							<List>
								<Subheader>
									{this.state.name}
									<IconButton tooltip="Edit" style={{float: "right"}}
										onTouchTap={() => this.editAccount.open()}>
										<FontIcon className="material-icons">mode_edit</FontIcon>
									</IconButton>
								</Subheader>
								<Divider />
								<ListItem
								  hoverColor="white"
								  primaryText="Address"
								  secondaryText={window.entry.address.encoded}/>
								<Divider />
								<ListItem
								  hoverColor="white"
								  primaryText="Description"
								  secondaryText={window.entry.desc}
								  secondaryTextLines={2}/>
								<Divider />
								<ListItem
								  hoverColor="white"
								  primaryText={<div> {window.entry.postCount} Posts</div>}/>
							</List>
						</div>
					</ToolTip>

					<Upload nav={this} />
					<EditAccount nav={this} />
				</ul>

			  )
	}
}

export default NavBar;
