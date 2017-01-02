import React, { Component } from 'react';
import { Router, Route, hashHistory, Link } from 'react-router';
import request from "superagent"

import {List, ListItem} from 'material-ui/List';
import Divider from 'material-ui/Divider';
import Subheader from 'material-ui/Subheader';

import util from "../util.js"
import ReactTooltip from 'react-tooltip'
import ToolTip from 'react-portal-tooltip'

class NavBar extends Component{

	constructor(props){
		super(props);


		this.state = { name: window.entry.name, uploadTooltip: false, channelTooltip: false };
		window.navbar = this;
	}

	toggleUploadTooltip(){
		this.setState({ uploadTooltip: !this.state.uploadTooltip, channelTooltip: false })
	}

	toggleChannelTooltip(){
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
								onMouseLeave={this.toggleUploadTooltip.bind(this)}>
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
								<Subheader>{window.entry.name}</Subheader>
								<Divider inset={true} />
								<ListItem
								  hoverColor="white"
								  primaryText="Address"
								  secondaryText={window.entry.address.encoded}/>
								<Divider inset={true} />
								<ListItem
								  hoverColor="white"
								  primaryText="Description"
								  secondaryText={window.entry.desc}
								  secondaryTextLines={2}/>
								<Divider inset={true} />
								<ListItem
								  hoverColor="white"
								  primaryText={<div> {window.entry.postCount} Posts</div>}/>
							</List>
						</div>
					</ToolTip>
				</ul>

			  )
	}
}

export default NavBar;
