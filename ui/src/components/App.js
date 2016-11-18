import React, { Component } from 'react';
import { Router, Route, hashHistory, Link } from 'react-router';

import AppBar from 'material-ui/AppBar';
import Drawer from 'material-ui/Drawer';
import MenuItem from 'material-ui/MenuItem';
import TextField from 'material-ui/TextField';
import {grey100, grey50} from 'material-ui/styles/colors';

import Home from './Home';
import Search from "./Search"
import SearchResults from "./SearchResults"
import Stream from "./Stream"

import util from "../util"

var WebTorrent = require("webtorrent");

var routes = [{ path: "/", component: Home },
			  { path: "/search", component: SearchResults },
			  { path: "/stream/:infohash", component: Stream },
			  { path: "/downloads", component: Stream }];

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

		this.loadConfig();

		this.handleToggle = this.handleToggle.bind(this);
		this.onResults = this.onResults.bind(this);

		this.config = util.loadConfig();

		window.downloadClient = new WebTorrent();
	}

	loadConfig()
	{ 
	}

	handleToggle(){ this.setState({ drawerOpen: !this.state.drawerOpen }) }

	onResults(res) 
	{
		routes[1].component = () => {
			return (
				<SearchResults posts={res}/>
			)
		};

		hashHistory.push("/search")
	}

	homeButtonClick(){ hashHistory.push("/") }

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

					<Search
						onResults={this.onResults}
						Subscriptions={this.config.subscriptions}/>

				</AppBar>

				<Drawer width={200} 
						docked={true} 
						open={this.state.drawerOpen} 
						containerClassName="drawer"
						containerStyle={style.drawer}>

					<div style={style.drawerItems}>
						<a onClick={this.homeButtonClick}><MenuItem>Home</MenuItem></a>
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
