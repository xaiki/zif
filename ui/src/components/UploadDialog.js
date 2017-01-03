import React, { Component } from 'react';
import request from "superagent";

import Dialog from 'material-ui/Dialog';
import FlatButton from 'material-ui/FlatButton';
import Subheader from 'material-ui/Subheader';
import Checkbox from 'material-ui/Checkbox';
import {List, ListItem} from 'material-ui/List';
import TextField from 'material-ui/TextField';
import Divider from 'material-ui/Divider';

import util from "../util"
import nt from "nt"
import fs from "fs"

class Upload extends Component
{
	constructor(props)
	{
		super(props);

		this.props.nav.upload = this;

		this.state = { open: false };

		this.handleClose = this.handleClose.bind(this);
		this.handleTitleChange = this.handleTitleChange.bind(this);
		this.handleDescChange = this.handleDescChange.bind(this);
		this.open = this.open.bind(this);

		this.title = this.desc = "";
	}

	open() {
		console.log("open")
		this.setState({ open: true });
	}

	handleClose() { 
		if (this.title.length == 0) {
			return;
		}

		var dirRegex = "^(.+)/([^/]+)$";

		var dir = this.uploadFile.files[0].path.match(dirRegex)[1];
		var files = [];

		for (var i = 0; i < this.uploadFile.files.length; i++) {
			files.push(this.uploadFile.files[i].name);
		}

		var rs = nt.make('udp://tracker.zer0day.to:1337/announce', 
					dir,
					files);
		rs.pipe(fs.createWriteStream(this.title + '.torrent'));

		// actually send the upload request
		var data = {
			title: this.title,
			meta: JSON.stringify({ description: this.desc })
		};

		request.post("http://127.0.0.1:8080/self/addpost/")
			.type("form")
			.send({data: JSON.stringify(data), index: true})
			.end((err, res) => {
				console.log(res);
				this.setState({ open: false });
			});
	}

	handleTitleChange(e) {
		this.title = e.target.value;
	}

	handleDescChange(e) {
		this.desc = e.target.value;
	}

	render() 
	{
		const actions = [
		  <FlatButton
			label="Cancel"
			primary={true}
			onTouchTap={()=> this.setState({ open: false})}
		  />,
		  <FlatButton
			label="Upload"
			primary={true}
			onTouchTap={this.handleClose}
		  />,
		];

		return (<Dialog
		  title="Upload"
		  modal={false}
		  open={this.state.open}
		  onRequestClose={this.handleClose}
		  actions={actions}>

			<input ref={(i) => this.uploadFile = i} type="file" multiple></input>
			<TextField
				onChange={this.handleTitleChange}
				floatingLabelText="Title"
				fullWidth={true}/><br/>
			<TextField
				onChange={this.handleDescChange}
				floatingLabelText="Description"
				fullWidth={true}
				multiLine={true}/><br/>


		</Dialog>)
	}
}

export default Upload;
