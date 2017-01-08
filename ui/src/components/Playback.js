import React, { Component } from 'react';

import Dialog from 'material-ui/Dialog';
import RaisedButton from 'material-ui/RaisedButton';
import Subheader from 'material-ui/Subheader';
import {List, ListItem} from 'material-ui/List';

import Video from 'react-html5video';

import ReactList from 'react-list';
import {Wave} from "better-react-spinkit";

import File from "./File";
import util from "../util"

var torrent = remote.require("torrent-stream");

class Playback extends Component
{
	constructor(props) 
	{
		super(props);

		this.state = {
			open: true
		};
	}

	componentDidMount(){
	
	}

	render() {
		return (<Dialog
		  modal={false}
		  open={this.state.open}
		  onRequestClose={() => {this.setState({ open: false}); this.props.onClose();}}>
		
		  <Video width="100%" controls autoPlay >
		  	<source src={this.props.url}></source>
		  </Video>
		
		</Dialog>)
	}
}

export default Playback;
