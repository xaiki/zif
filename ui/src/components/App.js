import '../assets/stylesheets/base.scss';
import '../assets/stylesheets/app.scss';

import React, { Component } from 'react';
import { Router, Route, hashHistory } from 'react-router'

import AppBar from 'material-ui/AppBar';
import Drawer from 'material-ui/Drawer';
import MenuItem from 'material-ui/MenuItem';

import Home from './Home';


class App extends Component
{

	constructor(props)
	{
		super(props);
		this.state = {
			drawerOpen : true
		};

		this.handleToggle = this.handleToggle.bind(this);
	}

	handleToggle(){ this.setState({ drawerOpen: !this.state.drawerOpen }) }

	render() 
	{
		return(
			<div>
				<AppBar title="Zif"
					onLeftIconButtonTouchTap={this.handleToggle}
				/>

				<Drawer width={200} docked={true} open={this.state.drawerOpen} className="drawer">
					<div className="drawerItems">
						<MenuItem>Home</MenuItem>
					</div>
				</Drawer>

				<div id="router">
					<Router history={hashHistory}>
						<Route path="/" component={Home}/>
					</Router>
				</div>
			</div>
		)
	}
}

export default App;
