
import React, { Component } from 'react';
import request from "superagent";

import Dialog from 'material-ui/Dialog';
import FlatButton from 'material-ui/FlatButton';
import Subheader from 'material-ui/Subheader';
import Checkbox from 'material-ui/Checkbox';
import {List, ListItem} from 'material-ui/List';
import TextField from 'material-ui/TextField';

import util from "../util"

class EditAccount extends Component
{

	constructor(props)
	{
		super(props);

		this.state = {
			open: false
		};

		this.name = "";
		this.desc = "";

		this.props.nav.editAccount = this;

		this.handleClose = this.handleClose.bind(this);
		this.handleNameChange = this.handleNameChange.bind(this);
		this.handleDescChange = this.handleDescChange.bind(this);

		this.setChannelDetails = this.setChannelDetails.bind(this);
	}

	handleClose(){
		this.setState({ open: false });

		this.setChannelDetails();
	}

	open(){
		this.setState({ open: true });
	}

	handleNameChange(event){
		this.name = event.target.value;
	}

	handleDescChange(event){
		this.desc = event.target.value;
	}

	setChannelDetails(){

		if (this.name.length > 0) {
			request.post("http://127.0.0.1:8080/self/set/name/")
					.type("form")
					.send({ value: this.name })
					.end((err, res) => {
					});

			this.props.nav.setState({ name: this.name });
			window.entry.name = this.name;
		}

		if (this.desc.length > 0) {
			request.post("http://127.0.0.1:8080/self/set/desc/")
					.type("form")
					.send({ value: this.desc})
					.end((err, res) => {
					});
			window.entry.desc = this.desc;
		}

	}

	render() {
		return (<Dialog
		  title="Edit account info"
		  modal={false}
		  open={this.state.open}
		  onRequestClose={this.handleClose}
		  actions={<FlatButton label="Save"
        						primary={true}
        						onTouchTap={this.handleClose}/>}>

			<p>Enter the details of your channel</p>
			<TextField
				onChange={this.handleNameChange}
				floatingLabelText="Name"
				fullWidth={true}
				defaultValue={window.entry.name}/><br/>

			<TextField
				onChange={this.handleDescChange}
				floatingLabelText="Description"
				multiLine={true}
				fullWidth={true}
				defaultValue={window.entry.desc}
				/><br/>


		</Dialog>)
	}
}

export default EditAccount;
