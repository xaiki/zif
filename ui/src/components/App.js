import styles from  '../assets/stylesheets/base.scss';

import React, { Component } from 'react';
import { Router, Route, hashHistory } from 'react-router';

import AppBar from 'material-ui/AppBar';
import Drawer from 'material-ui/Drawer';
import MenuItem from 'material-ui/MenuItem';
import TextField from 'material-ui/TextField';
import {grey100, grey50} from 'material-ui/styles/colors';

import Home from './Home';
import Search from "./Search"

const routes = [{path: "/", component: Home}];

class App extends Component
{

	constructor(props)
	{
		super(props);

		this.state = {
			drawerOpen : true,

			search: {
				focus: false
			}
		};

		this.handleToggle = this.handleToggle.bind(this);
	}

	handleToggle(){ this.setState({ drawerOpen: !this.state.drawerOpen }) }

	render() 
	{
		var style = {
			drawer: {
				backgroundColor: grey100,
				zIndex: 1000
			},

			drawerItems: {
				marginTop: "75px"
			},

			router: {
				paddingLeft: "210px",
				paddingRight: "210px",
				paddingBottom: "10px",
				marginTop: "75px"
			}

		}

		return(
			<div style={{height: "100%"}}>
				<AppBar 
					title="Zif"
					style={{position: "fixed", top: 0, paddingRight: 0}}
					onLeftIconButtonTouchTap={this.handleToggle}>

					<Search/>

				</AppBar>

				<Drawer width={200} 
						docked={true} 
						open={this.state.drawerOpen} 
						containerClassName="drawer"
						containerStyle={style.drawer}>

					<div style={style.drawerItems}>
						<MenuItem>Home</MenuItem>
					</div>

				</Drawer>

				<div style={style.router}>
					<Router history={hashHistory} routes={routes}>
					</Router>
				</div>
			</div>
		)
	}
}

export default App;
