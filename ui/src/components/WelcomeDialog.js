import React, { Component } from 'react';
import request from "request";

import Dialog from 'material-ui/Dialog';
import FlatButton from 'material-ui/FlatButton';
import Subheader from 'material-ui/Subheader';
import Checkbox from 'material-ui/Checkbox';
import {List, ListItem} from 'material-ui/List';

import util from "../util"

class Welcome extends Component
{
	constructor(props)
	{
		super(props);

		this.config = this.props.config;

		this.state = {
			open: !this.config.welcomed,
			bootstrap: {
				zif: false
			}
		};

		this.handleClose = this.handleClose.bind(this);

		this.checkZif = this.checkZif.bind(this);
	}

	handleClose() { 
		// Make sure at least one has been selected.
		if (!this.state.bootstrap.zif) return;

		// BOOTSTRAP! :D
		if (this.state.bootstrap.zif) {
			this.bootstrapUrl("zif.io")
		}

		this.setState({ open: !this.state.open });

		this.config.welcomed = true;
		util.saveConfig(this.config);
	}

	bootstrapUrl(url) {
		request('http://127.0.0.1:8080/self/bootstrap/' + url + "/",
					function (error, response, body){ console.log(response); })
	}

	checkZif(e, checked) {
		this.setState({
			bootstrap: { zif: checked }
		});
	}

	render() 
	{
		return (<Dialog
		  title="Welcome to Zif"
		  modal={false}
		  open={this.state.open}
		  onRequestClose={this.handleClose}
		  actions={<FlatButton label="Go"
        						primary={true}
        						onTouchTap={this.handleClose}
        						disabled={!(this.state.bootstrap.zif)}/>}
		>
			<p>Pick one or more of the channels below to get started</p>
			<List>
				  <ListItem 
				  		primaryText="zif.io" 
				  		leftCheckbox={<Checkbox onCheck={this.checkZif} />}
				  		secondaryText={
				  			<p>The official Zif channel</p>
				  		}/>
    		</List>,

		</Dialog>)
	}
}

export default Welcome;
