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
var wjs = require("wcjs-player");


class Playback extends Component
{
	constructor(props) 
	{
		super(props);

		this.state = {
			open: true
		};

		this.componentDidMount = this.componentDidMount.bind(this);
	}

	componentDidMount(){
		this.player = new wjs("#player").addPlayer({
			  autoplay: true,
			   wcjs: require('webchimera.js')
			});

			this.player.addPlaylist(this.props.url);

			if (this.playerDOM.style.height != this.playerDOM.parentNode.style.maxHeight){
				this.playerDOM.style.height = this.playerDOM.parentNode.style.maxHeight;
				this.playerDOM.parentNode.style.padding = "0";
				this.forceUpdate();
			}
	}

	render() {
		return (<Dialog
		  modal={false}
		  open={this.state.open}
		  onRequestClose={() => {this.setState({ open: false}); this.props.onClose();}}
		  style={{padding: "0"}}>
		
		  <div ref={(i) => this.playerDOM = i } id="player" 
		  		style={{width: "100%"}}></div>
		
		</Dialog>)
	}
}

export default Playback;
