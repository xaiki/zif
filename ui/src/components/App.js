import styles from  '../assets/stylesheets/base.scss';

import React, { Component } from 'react';
import { Router, Route, hashHistory } from 'react-router';

import AppBar from 'material-ui/AppBar';
import Drawer from 'material-ui/Drawer';
import MenuItem from 'material-ui/MenuItem';
import TextField from 'material-ui/TextField';
import {grey100, grey50} from 'material-ui/styles/colors';



import Home from './Home';

const style = {
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
	},

	search: {
		backgroundColor: "white",
		paddingLeft: "5px",
		paddingRight: "5px",
		width: "100%"
	}

}

class App extends Component
{

	constructor(props)
	{
		super(props);

		this.state = {
			drawerOpen : true,

			search: {
				width: "256px",
				height: "48px",
				postition: "fixed",
				top: "0",
				left: "0"
			}
		};

		this.handleToggle = this.handleToggle.bind(this);
	}

	handleToggle(){ this.setState({ drawerOpen: !this.state.drawerOpen }) }

	handleSearchFocus() 
	{
		this.setState({ search: 
			{ 
				width: "400px", 
			}});
	}

	render() 
	{
		return(
			<div style={{height: "100%"}}>
				<AppBar 
					title="Zif"
					style={{position: "fixed", top: 0}}
					onLeftIconButtonTouchTap={this.handleToggle}>

					<div style={this.state.search}>
						<TextField 
							style={style.search} 
							underlineShow={false}
							hintText="Search..."
						/>
					</div>
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
					<Router history={hashHistory}>
						<Route path="/" component={Home}/>
					</Router>
				</div>
			</div>
		)
	}
}

export default App;
